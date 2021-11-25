package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

/*
	About sheet.
*/
const (
	defautWriteRowCount = 5000
)

func newSheet(id string, client *SpreadSheets) *Sheet {
	return &Sheet{
		ssClient: client,
		id:       id,
	}
}

// Sheet represent a sheet tab in spread sheets(SpreadSheets)
type Sheet struct {
	Err      error
	ssClient *SpreadSheets
	id       string // sheet id
}

func (s *Sheet) GetID() string {
	return s.id
}

func (s *Sheet) GetContentByRange(startCellname, endCellname string) (*SheetContent, error) {
	r := s.genRange(startCellname, endCellname)
	url := fmt.Sprintf("%s/open-apis/sheet/v2/spreadsheets/%s/values/%s?valueRenderOption=ToString", s.ssClient.baseClient.domain, s.ssClient.token, r)
	_req, _ := http.NewRequest(http.MethodGet, url, nil)
	var content *SheetContent
	_, err := s.ssClient.baseClient.CommonReq(_req, &content)
	return content, err
}

// GetContentByRangeV2
// Reference https://open.feishu.cn/document/ukTMukTMukTM/ugTMzUjL4EzM14COxMTN
func (s *Sheet) GetContentByRangeV2(startCellname, endCellname string, render SheetRenderOption, dateTime SheetDateTimeRenderOption) (*SheetContent, error) {
	r := s.genRange(startCellname, endCellname)
	u := fmt.Sprintf("%s/open-apis/sheet/v2/spreadsheets/%s/values/%s", s.ssClient.baseClient.domain, s.ssClient.token, r)
	params := url.Values{}
	if render != "" {
		params.Set("valueRenderOption", string(render))
	}
	if dateTime != "" {
		params.Set("dateTimeRenderOption", string(dateTime))
	}
	if len(params) > 0 {
		u = u + "?" + params.Encode()
	}

	_req, _ := http.NewRequest(http.MethodGet, u, nil)
	var content *SheetContent
	_, err := s.ssClient.baseClient.CommonReq(_req, &content)
	return content, err
}

// GetRows for get all rows
func (s *Sheet) GetRows(withFirstLine bool) ([]SheetRow, error) {
	meta, err := s.getMeta()
	if err != nil {
		return nil, err
	}
	content, err := s.GetContentByRangeV2("A1", fmt.Sprintf("%s%d", num2ColName(meta.ColumnCount), meta.RowCount), SheetRenderFormattedValue, SheetDateTimeRenderFormattedString)
	if err != nil {
		return nil, err
	}
	rows := content.ToRows()
	if len(rows) > 1 && !withFirstLine {
		rows = rows[1:]
	}
	return rows, nil
}

// WriteRows write rows line by line, start from A1 cell
func (s *Sheet) WriteRows(title []string, data [][]interface{}, batchCount ...int) *Sheet {
	return s.WriteRowsByStartCell("A1", title, data, batchCount...)
}

// WriteRowsByStartCell
// Parameter
//  title: title of every columns.
//  batchCount: max insert line coune once.
// Example
//  s.WriteRowsByStartCell("A1",[]string{"name", "age"}, [][]interface{}{
//  	{"Ace",15},
//  	{"Bob",16},
//  },10)
func (s *Sheet) WriteRowsByStartCell(startCell string, title []string, data [][]interface{}, batchCount ...int) *Sheet {
	if s.Err != nil {
		return s
	}
	hasTitile := len(title) > 0
	hasData := len(data) > 0
	if !hasTitile && !hasData {
		s.Err = newErr("title and data is all empty")
		return s
	}
	if hasData && len(data[0]) > 100 {
		s.Err = newErr("max column count is 100, yous is %d", len(data[0]))
		return s
	}
	if hasTitile {
		newData := make([][]interface{}, 0, len(data)+1)
		titles := make([]interface{}, 0, len(title))
		for _, v := range title {
			titles = append(titles, v)
		}
		newData = append(newData, titles)
		newData = append(newData, data...)
		data = newData
	}
	u := fmt.Sprintf("%s/open-apis/sheets/v2/spreadsheets/%s/values", s.ssClient.baseClient.domain, s.ssClient.token)
	colCount := len(data[0])
	//onceWriteRowCount := 1000 / colCount // 一次写入的行数
	rowCount := len(data)
	onceWriteRowCount := defautWriteRowCount
	if len(batchCount) > 0 {
		onceWriteRowCount = batchCount[0]
	}
	//if onceWriteRowCount == 0 {
	//	onceWriteRowCount = 1
	//}
	writeTimes := rowCount / onceWriteRowCount
	if rowCount%onceWriteRowCount > 0 {
		writeTimes += 1
	}
	for i := 0; i < writeTimes; i++ {
		startRow := i * onceWriteRowCount
		endRow := (i + 1) * onceWriteRowCount
		if endRow > rowCount {
			endRow = rowCount
		}
		writeRowCount := endRow - startRow
		endCell := cellnameAdd(startCell, colCount-1, writeRowCount-1) // 因为 start cell 所在的行也可以写入，所以减 1
		body, _ := json.Marshal(map[string]interface{}{
			"valueRange": map[string]interface{}{
				"range":  s.genRange(startCell, endCell),
				"values": data[startRow:endRow],
			},
		})
		_req, _ := http.NewRequest(http.MethodPut, u, bytes.NewReader(body))
		_, err := s.ssClient.baseClient.CommonReq(_req, nil)
		if err != nil {
			s.Err = newErr(err.Error())
			return s
		}
		startCell = cellnameAdd(startCell, 0, writeRowCount)
	}
	return s
}

// UpdateTitle ...
func (s *Sheet) UpdateTitle(title string) *Sheet {
	return s.updateBase(map[string]interface{}{
		"title": title,
	})
}

// UpdateIndex ...
func (s *Sheet) UpdateIndex(index int) *Sheet {
	return s.updateBase(map[string]interface{}{
		"index": index,
	})
}

func (s *Sheet) Hidden(hidden bool) *Sheet {
	return s.updateBase(map[string]interface{}{
		"hidden": hidden,
	})
}

// FrozenRow
// Parameter
//  row: number of row that want to frezen. 0 represent unfrozen
func (s *Sheet) FrozenRow(row int) *Sheet {
	return s.updateBase(map[string]interface{}{
		"frozenRowCount": row,
	})
}

// FrozenColumn
// Parameter
//  column: number of column that want to frezen. 0 represent unfrozen
func (s *Sheet) FrozenColumn(column int) *Sheet {
	return s.updateBase(map[string]interface{}{
		"frozenColCount": column,
	})
}

func (s *Sheet) Protect(info string, userIDs []string) *Sheet {
	m := map[string]interface{}{
		"lock":     "LOCK",
		"lockInfo": info,
	}
	if len(userIDs) > 0 {
		m["userIDs"] = userIDs
	}
	return s.updateBase(m)
}

func (s *Sheet) MoveRows(start, end, target int) *Sheet {
	return s.moveRowsOrCols(start, end, target, true)
}

func (s *Sheet) MoveColumns(start, end, target int) *Sheet {
	return s.moveRowsOrCols(start, end, target, false)
}

func (s *Sheet) moveRowsOrCols(start, end, target int, row bool) *Sheet {
	if s.Err != nil {
		return s
	}
	flag := "ROWS"
	if !row {
		flag = "COLUMNS"
	}
	_url := s.ssClient.baseClient.urlJoin(fmt.Sprintf("/open-apis/sheets/v3/spreadsheets/%s/sheets/%s/move_dimension", s.ssClient.token, s.id))
	en, _ := json.Marshal(map[string]interface{}{
		"source": map[string]interface{}{
			"major_dimension": flag,
			"start_index":     start,
			"end_index":       end,
		},
		"destination_index": target,
	})
	req, _ := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(en))
	_, err := s.ssClient.baseClient.CommonReq(req, nil)
	s.Err = err
	return s
}

func (s *Sheet) updateBase(m map[string]interface{}) *Sheet {
	if s.Err != nil {
		return s
	}
	m["sheetId"] = s.id
	_, _, err := s.ssClient.origin.SheetBatchUpdate(map[ModifySheetType]interface{}{
		ModifySheetUpdate: map[string]interface{}{
			"properties": m,
		},
	})
	s.Err = err
	return s
}

func (s *Sheet) getMeta() (*sheetMeta, error) {
	meta, err := s.ssClient.GetMeta()
	if err != nil {
		return nil, fmt.Errorf("get sheet meta, sheetID: %s, err: %s", s.id, err.Error())
	}
	for _, v := range meta.Sheets {
		if v.SheetID == s.id {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("get sheet meta, sheetID: %s, can not find this sheet", s.id)
}

// SheetRow represent for a group of sheet cell
type SheetRow = []*SheetCell
