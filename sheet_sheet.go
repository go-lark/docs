package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

/*
About sheet.
*/
const (
	defautWriteRowCount = 5000
)

func newSheet(id string, client *SpreadSheets) *Sheet {
	s := &Sheet{
		ssClient: client,
		id:       id,
	}
	//s.SheetRange = &SheetRange{sheet: s, rangeVal: id}
	return s
}

// Sheet represent a sheet tab in spread sheets(SpreadSheets)
type Sheet struct {
	Err      error
	ssClient *SpreadSheets
	id       string // sheet id

	//	*SheetRange
}

func (s *Sheet) GetID() string {
	return s.id
}

func (s *Sheet) GetContentByRange(startCellname, endCellname string) (*SheetContent, error) {
	r := s.genRange(startCellname, endCellname)
	return s.getContentByRange(r)
}

func (s *Sheet) getContentByRange(r string) (*SheetContent, error) {
	url := fmt.Sprintf("%s/open-apis/sheet/v2/spreadsheets/%s/values/%s", s.ssClient.baseClient.domain, s.ssClient.token, r)
	_req, _ := http.NewRequest(http.MethodGet, url, nil)
	var content *SheetContent
	_, err := s.ssClient.baseClient.CommonReq(_req, &content)
	return content, err
}

// !A1:D5
func (s *Sheet) NewRangeFull(startCellName, endCellName string) *SheetRange {
	r := &SheetRange{}
	if s.Err != nil {
		r.Err = s.Err
		return r
	}
	r.leftTop = NewCellName(startCellName)
	r.rightBottom = NewCellName(endCellName)
	r.sheet = s
	return r
}

/*
// !A1:D
func (s *Sheet) NewRangeHalf(startCellName, endCol string) *SheetRange {
	r := &SheetRange{}
	if s.Err != nil {
		r.Err = s.Err
		return r
	}
	r.rangeVal = s.id + "!" + startCellName + ":" + endCol
	r.sheet = s
	return r
}

// !A:D
func (s *Sheet) NewRangeCol(startCol, endCol string) *SheetRange {
	r := &SheetRange{}
	if s.Err != nil {
		r.Err = s.Err
		return r
	}
	r.rangeVal = s.id + "!" + startCol + ":" + endCol
	r.sheet = s
	return r
}

func (s *Sheet) NewRangeRow(startRow, endRow string) *SheetRange {
	r := &SheetRange{}
	if s.Err != nil {
		r.Err = s.Err
		return r
	}
	meta, err := s.getMeta()
	if err != nil {
		r.Err = err
		return r
	}
}
*/

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
	content := &SheetContent{}
	_, err := s.ssClient.baseClient.CommonReq(_req, &content)
	content.sheet = s
	return content, err
}

// ReadRows for get all rows
func (s *Sheet) ReadRows() ([]SheetRow, error) {
	meta, err := s.getMeta()
	if err != nil {
		return nil, err
	}
	content, err := s.GetContentByRangeV2("A1", fmt.Sprintf("%s%d", num2ColName(meta.ColumnCount), meta.RowCount), SheetRenderFormattedValue, SheetDateTimeRenderFormattedString)
	if err != nil {
		return nil, err
	}
	return content.ToRows()
}

func (s *Sheet) TrimBlankTail(rows []SheetRow) []SheetRow {
	w := len(rows) - 1
	if w == 0 {
		return rows
	}
	for {
		for _, v := range rows[w] {
			if v.val != nil {
				return rows[:w+1]
			}
		}
		w--
		if w == 0 {
			return rows[:w]
		}
	}
}

// WriteRows write rows line by line, start from A1 cell
func (s *Sheet) WriteRows(data [][]interface{}, batchCount ...int) *Sheet {
	return s.WriteRowsByStartCell("A1", nil, data, batchCount...)
}

// WriteRowsByStartCell
// Parameter
//
//	title: title of every columns.
//	batchCount: max insert line coune once.
//
// Example
//
//	s.WriteRowsByStartCell("A1",[]string{"name", "age"}, [][]interface{}{
//		{"Ace",15},
//		{"Bob",16},
//	},10)
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
//
//	row: number of row that want to frezen. 0 represent unfrozen
func (s *Sheet) FrozenRow(row int) *Sheet {
	return s.updateBase(map[string]interface{}{
		"frozenRowCount": row,
	})
}

// FrozenColumn
// Parameter
//
//	column: number of column that want to frezen. 0 represent unfrozen
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

// MoveRows start with 1
func (s *Sheet) MoveRows(start, end, target int) *Sheet {
	return s.moveRowsOrCols(start, end, target, true)
}

// MoveColumns start with 1
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

/*
func (s *Sheet) getRange() *SheetRange {
	r := &SheetRange{}
	if s.Err != nil {
		r.Err = s.Err
		return r
	}
	meta, err := s.getMeta()
	if err != nil {
		r.Err = s.Err
		return r
	}
}
*/

// SheetRow represent for a group of sheet cell
type SheetRow = []*SheetCell

/*
type mergeInfo struct {
	startCol int
	endCol   int
	startRow int
	endRow   int
}
*/

type (
	SheetContent struct {
		sheet      *Sheet
		ValueRange valueRange `json:"valueRange"`
		Err        error
	}
	valueRange struct {
		Values [][]interface{} `json:"values"`
	}
)

func (sc *SheetContent) ToRows() ([]SheetRow, error) {
	if sc.Err != nil {
		return nil, sc.Err
	}
	sheetRows := make([]SheetRow, 0)
	if sc == nil {
		return sheetRows, nil
	}
	for _, rows := range sc.ValueRange.Values {
		cells := make([]*SheetCell, 0)
		for _, row := range rows {
			r := row
			cells = append(cells, NewSheetCell(r))
		}
		sheetRows = append(sheetRows, cells)
	}
	return sheetRows, nil
}

func (sc *SheetContent) ToRowsTrimBlankTail() ([]SheetRow, error) {
	rows, err := sc.ToRows()
	if err != nil {
		return nil, err
	}
	w := len(rows) - 1
	if w == 0 {
		return rows, nil
	}
	for {
		for _, v := range rows[w] {
			if v.val != nil {
				return rows[:w+1], nil
			}
		}
		w--
		if w == 0 {
			return rows[:w+1], nil
		}
	}
}

// RowsParseMerge
// parse cell from merged cell.
// if a cell is merged, like below, we only get value at the first cell, others are nil.
// then we fill a to every cell of the merged cell.
// | a |   |   |
// |   |   |   |
// |   |   |   |
func (sc *SheetContent) ToRowsParseMerged() ([]SheetRow, error) {
	meta, err := sc.sheet.getMeta()
	if err != nil {
		return nil, err
	}
	cells, err := sc.ToRows()
	if err != nil {
		return nil, err
	}
	for _, v := range meta.Merges {
		startCol := v.StartColumnIndex
		endCol := v.StartColumnIndex + v.ColumnCount - 1
		startRow := v.StartRowIndex
		endRow := v.StartRowIndex + v.RowCount - 1
		for i := startRow; i <= endRow; i++ {
			for j := startCol; j <= endCol; j++ {
				cells[i][j] = cells[startRow][startCol]
			}
		}
	}
	return cells, nil
}

func (sc *SheetContent) Scan(ptr interface{}) error {
	rows, err := sc.ToRowsParseMerged()
	if err != nil {
		return err
	}
	return sc.scan(rows, ptr)
}

func (sc *SheetContent) scan(rows []SheetRow, ptr interface{}) error {
	// check it args is a pointer
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("ptr is must a slice pointer")
	}
	if rv.CanSet() {
		return fmt.Errorf("can not get slice address")
	}
	// dereference the pointer to get the slice
	rv = rv.Elem()
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("ptr element is must a slice")
	}
	// slice element
	elemt := rv.Type().Elem()
	if elemt.Kind() != reflect.Ptr {
		return fmt.Errorf("slice element must be a struct pointer, it is not pointer")
	}
	elemt = elemt.Elem()
	if elemt.Kind() != reflect.Struct {
		return fmt.Errorf("slice element must be a struct pointer, it it not struct")
	}

	fieldPosition := map[int]int{}
	for i := 0; i < elemt.NumField(); i++ {
		str := elemt.Field(i).Tag.Get("docs")
		if str != "" {
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return fmt.Errorf("struct tag wrong, it is not a number,  field index, %d, tag: %s", i, str)
			}
			fieldPosition[i] = int(v)
		} else {
			fieldPosition[i] = i
		}
	}
	for _, row := range rows {
		// new a slice element
		valP := reflect.New(elemt)
		// dereference pointer
		val := valP.Elem()
		for j := 0; j < val.NumField(); j++ {
			name := elemt.Field(j).Name
			position := fieldPosition[j]
			f := val.Field(j)
			if !f.CanAddr() {
				fmt.Printf(": can not set, %s\n", name)
			}
			switch f.Kind() {
			case reflect.Int64, reflect.Int:
				to, err := row[position].ToInt64()
				if err != nil {
					return fmt.Errorf("scan faild, index: %d, %w", j, err)
				}
				f.SetInt(to)
			case reflect.Float64, reflect.Float32:
				to, err := row[position].ToFloat()
				if err != nil {
					return fmt.Errorf("scan faild, index: %d, %w", j, err)
				}
				f.SetFloat(to)
			default:
				f.SetString(row[position].ToString())
			}

		}
		rv.Set(reflect.Append(rv, valP))
	}
	return nil
}
