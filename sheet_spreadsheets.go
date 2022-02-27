package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func newSpreadSheet(token string, client *Client) *SpreadSheet {
	s := &SpreadSheet{
		baseClient: client,
	}
	s.token = token
	return s
}

// SpreadSheet represent a group of sheet
type SpreadSheet struct {
	Err error
	tokenIns
	baseClient *Client
	origin     *SpreadSheetOrigin
}

// GetOrigin get origin client, with is use origin open API.
// 获取原始客户端，这个客户端直接使用开放平台文档上的 API，没有做任何封装
func (ss *SpreadSheet) GetOrigin() *SpreadSheetOrigin {
	if ss.origin == nil {
		ss.origin = newSpreadSheetOrigin(ss.baseClient, ss.token)
	}
	return ss.origin
}

// GetMeta get spread sheet meta information.
// 获取这个表格的元信息。
func (ss *SpreadSheet) GetMeta() (res *SpreadSheetMeta, err error) {
	if ss.Err != nil {
		return nil, ss.Err
	}
	_, res, err = ss.GetOrigin().MetaInfo()
	return res, err
}

// UpdateTitle
// Parameter
//  title: tile of spreadsheet
func (ss *SpreadSheet) UpdateTitle(title string) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	_, err := ss.GetOrigin().Properties(&SpreadSheetProperties{
		Title: &title,
	})
	ss.Err = err
	return ss
}

// DeleteSheet
// Note:
//  Sheet id can be found in url, for example in the url https://laily.feishu.cn/sheets/shtcnLML6348M7ujOaYd1EsUe9f?sheet=5d8cef
// sheet id is 5d8cef
func (ss *SpreadSheet) DeleteSheet(sheetID string) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	args := map[string]string{
		"sheetId": sheetID,
	}
	_, _, err := ss.GetOrigin().SheetBatchUpdate(map[ModifySheetType]interface{}{
		ModifySheetDelete: args,
	})
	ss.Err = err
	return ss
}

func (ss *SpreadSheet) DeleteSelf() error {
	return newFile(ss.baseClient).deleteSpreadSheet(ss.GetToken())
}

func (ss *SpreadSheet) ModifyProperties(args *ModifyProperties) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	_, _, err := ss.GetOrigin().SheetBatchUpdate(map[ModifySheetType]interface{}{
		ModifySheetUpdate: map[string]interface{}{
			"properties": args,
		},
	})
	ss.Err = newErr(err.Error())
	return ss
}

// Share this document to user or group.
func (ss *SpreadSheet) Share(perm Perm, notify bool, members ...*Member) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	err := ss.baseClient.permission().Add(ss.token, FileTypeSheet, perm, notify, members...)
	ss.Err = err
	return ss
}

// UnShare this document to other user or group.
func (ss *SpreadSheet) UnShare(members ...*Member) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	err := ss.baseClient.permission().Remove(ss.token, FileTypeSheet, members...)
	ss.Err = err
	return ss
}

func (ss *SpreadSheet) ChangeOwner(newOwner *Member, removeOldOwner, notify bool) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	_, _, err := ss.baseClient.permission().TransferOwner(ss.token, FileTypeSheet, newOwner, removeOldOwner, notify)
	if err != nil {
		ss.Err = newErr(err.Error())
	}
	return ss
}

func (ss *SpreadSheet) SetAccessPermission(per string) *SpreadSheet {
	return ss.setPublic(&PublicSet{
		LinkShareEntity: &per,
	})
}

func (ss *SpreadSheet) setPublic(args *PublicSet) *SpreadSheet {
	if ss.Err != nil {
		return ss
	}
	err := ss.baseClient.permission().PublicSet(args)
	ss.Err = newErr(err.Error())
	return ss
}

// CreateSheet
// Parameter
//  index: first position is 0
func (ss *SpreadSheet) CreateSheet(title string, index int) *Sheet {
	sheet := &Sheet{}
	if ss.Err != nil {
		sheet.Err = ss.Err
		return sheet
	}
	args := &ModifySheet{
		Properties: &ModifyProperties{
			Title: &title,
			Index: &index,
		},
	}
	b, res, err := ss.GetOrigin().SheetBatchUpdate(map[ModifySheetType]interface{}{
		ModifySheetAdd: args,
	})
	if err != nil {
		sheet.Err = err
		return sheet
	}
	if len(res.Replies) == 0 {
		sheet.Err = newErr("sheet batch update res is empty, res: %s", string(b))
		return sheet
	}
	id := res.Replies[0].AddSheet.Properties.SheetID
	return ss.GetSheetByID(id)
}

// SheetID get sheet by sheet id
// 根据 sheet id 获取 sheet 实例。
func (ss *SpreadSheet) GetSheetByID(sheetID string) *Sheet {
	meta, err := ss.GetMeta()
	s := &Sheet{}
	if err != nil {
		s.Err = newErr(err.Error())
		return s
	}
	for _, vv := range meta.Sheets {
		v := vv
		if v.SheetID == sheetID {
			return newSheet(&v, ss)
		}
	}
	s.Err = newErr("sheet id not exist, id: %s", sheetID)
	return s
}

// SheetIndex get a sheet instance by index. Index start from 1
func (ss *SpreadSheet) GetSheetByIndex(index int) *Sheet {
	meta, err := ss.GetMeta()
	s := &Sheet{}
	if err != nil {
		s.Err = newErr(err.Error())
		return s
	}
	for _, vv := range meta.Sheets {
		v := vv
		if v.Index == index {
			return newSheet(&v, ss)
		}
	}
	s.Err = newErr("sheet index not exist, index: %d", index)
	return s
}

// SheetName get sheet by sheet name
func (ss *SpreadSheet) GetSheetByName(name string) *Sheet {
	meta, err := ss.GetMeta()
	s := &Sheet{}
	if err != nil {
		s.Err = newErr(err.Error())
		return s
	}
	for _, vv := range meta.Sheets {
		v := vv
		if v.Title == name {
			return newSheet(&v, ss)
		}
	}
	s.Err = newErr("sheet name not exist, name: %s", name)
	return s
}

// CopySheet
// Pareameter
//  sourceSheetID: sheet ID of which will be copied.
//  title: title of new sheet.
func (ss *SpreadSheet) CopySheet(sourceSheetID string, title string) (sheet *Sheet) {
	sheet = &Sheet{}
	if ss.Err != nil {
		sheet.Err = ss.Err
		return
	}
	args := map[string]map[string]string{
		"source": {
			"sheetId": sourceSheetID,
		},
		"destination": {
			"title": title,
		},
	}
	b, res, err := ss.GetOrigin().SheetBatchUpdate(map[ModifySheetType]interface{}{
		ModifySheetCopy: args,
	})
	if err != nil {
		sheet.Err = newErr(err.Error())
		return
	}
	if len(res.Replies) == 0 {
		sheet.Err = newErr("sheet batch update res is empty, res: %s", string(b))
		return sheet
	}
	id := res.Replies[0].CopySheet.Properties.SheetID
	return ss.GetSheetByID(id)
}

type SheetContent struct {
	ValueRange struct {
		Values [][]interface{} `json:"values"`
	} `json:"valueRange"`
}

func (sc *SheetContent) ToRows() []SheetRow {
	sheetRows := make([]SheetRow, 0)
	if sc == nil {
		return sheetRows
	}
	for _, rows := range sc.ValueRange.Values {
		cells := make([]*SheetCell, 0)
		for _, row := range rows {
			r := row
			cells = append(cells, NewSheetCell(r))
		}
		sheetRows = append(sheetRows, cells)
	}
	return sheetRows
}

type (
	ModifySheet struct {
		Properties *ModifyProperties `json:"properties,omitempty"`
	}
	ModifyProperties struct {
		SheetID *string        `json:"sheetId,omitempty"`
		Title   *string        `json:"title,omitempty"`
		Index   *int           `json:"index,omitempty"`
		Hidden  *bool          `json:"hidden,omitempty"`
		Protect *ModifyProtect `json:"protect,omitempty"`
	}
	ModifyProtect struct {
		Lock     *string `json:"lock,omitempty"`
		LockInfo *string `json:"lockInfo,omitempty"`
		UserIds  []int64 `json:"userIds,omitempty"`
	}
)

type (
	// SpreadsheetMeta sheet 的 meta 信息
	SpreadSheetMeta = MetaInfoResp
	MetaInfoResp    struct {
		Properties       metaProp    `json:"properties"`
		Sheets           []sheetMeta `json:"sheets"`
		SpreadsheetToken string      `json:"spreadsheetToken"`
	}
	metaProp struct {
		Title       string `json:"title"`
		OwnerUser   int64  `json:"ownerUser"`
		OwnerUserID string `json:"ownerUserID"`
		SheetCount  int    `json:"sheetCount"`
		Revision    int    `json:"revision"`
	}
	sheetMeta struct {
		SheetID        string              `json:"sheetId"`
		Title          string              `json:"title"`
		Index          int                 `json:"index"`
		RowCount       int                 `json:"rowCount"`
		ColumnCount    int                 `json:"columnCount"`
		FrozenRowCount int                 `json:"frozenRowCount,omitempty"`
		FrozenColCount int                 `json:"frozenColCount,omitempty"`
		Merges         []*sheetMetaMerge   `json:"merges,omitempty"`
		Protect        []*sheetMetaProtect `json:"protectedRange"`
		BlockInfo      *sheetMetaBlock     `json:"blockInfo"`
	}

	sheetMetaMerge struct {
		ColumnCount      int `json:"columnCount"`
		RowCount         int `json:"rowCount"`
		StartColumnIndex int `json:"startColumnIndex"`
		StartRowIndex    int `json:"startRowIndex"`
	}

	sheetMetaProtect struct {
		Dimension sheetMetaProtectDimension `json:"dimension"`
		ProtectID string                    `json:"protectId"`
		SheetID   string                    `json:"sheetId"`
		LockInfo  string                    `json:"lockInfo"`
	}

	sheetMetaProtectDimension struct {
		EndIndex       int    `json:"endIndex"`
		MajorDimension string `json:"majorDimension"`
		SheetID        string `json:"sheetId"`
		StartIndex     int    `json:"startIndex"`
	}
	sheetMetaBlock struct {
		BlockToken string `json:"blockToken"`
		BlockType  string `json:"blockType"`
	}
)

type respBody struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// !A1:D5
func NewRangeFull(sheetID, startCellName, endCellName string) Range {
	return sheetID + "!" + startCellName + ":" + endCellName
}

// !A1:D
func NewRangeHalf(sheetID, startCellName, endCol string) Range {
	return sheetID + "!" + startCellName + ":" + endCol
}

// !A:D
func NewRangeCol(sheetID, startCol, endCol string) Range {
	return sheetID + "!" + startCol + ":" + endCol
}

// !<sheetID>
func NewRangeSheetID(sheetID string) Range {
	return sheetID
}

// Range reference https://open.feishu.cn/document/ukTMukTMukTM/uczNzUjL3czM14yN3MTN#bae19f77
type Range = string

func (s *Sheet) genRange(startCellname, endCellname string) string {
	// <sheetID>!A1:D5
	r := bytes.Buffer{}
	r.WriteString(s.GetID())
	r.WriteByte('!')
	r.WriteString(startCellname)
	r.WriteByte(':')
	r.WriteString(endCellname)
	return r.String()
}

func cellnameAdd(cellname string, colCount, rowCount int) string {
	cellname = strings.ToUpper(cellname)
	colname, rowname := cellnameSplit(cellname)
	return fmt.Sprintf("%s%d", colnameAdd(colname, colCount), rowname+rowCount)
}
