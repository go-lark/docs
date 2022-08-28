package docs

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

var (
	sheet *Sheet
)

/*
func TestReadContent(t *testing.T) {
	Convey("TestReadContent", t, func() {
		token := "shtcnPsIPv4woYKzSizh4rruCJe"
		sheet := getSheetClient().baseClient.OpenSpreadSheets(token).SheetIndex(1)
		So(sheet.Err, ShouldBeNil)
		rows, err := sheet.NewRangeFull("A1", "B2").Content().ToRows()

	})
}
*/

func TestSheetContent(t *testing.T) {
	Convey("TestSheetContent", t, func() {
		content := &SheetContent{
			ValueRange: valueRange{
				Values: [][]interface{}{
					{"ace", nil},
					{nil, nil},
				},
			},
		}

		Convey("not perse merge, not trim blank tail", func() {
			rows, err := content.ToRows()
			So(err, ShouldBeNil)
			So(rows, ShouldResemble, []SheetRow{
				{NewSheetCell("ace"), NewSheetCell(nil)},
				{NewSheetCell(nil), NewSheetCell(nil)},
			})
		})
		Convey("not perse merge, trim blank tail", func() {
			rows, err := content.ToRowsTrimBlankTail()
			So(err, ShouldBeNil)
			So(rows, ShouldResemble, []SheetRow{
				{NewSheetCell("ace"), NewSheetCell(nil)},
			})
		})
		Convey("perse merge, not trim blank tail", func() {
			var mSheet *Sheet
			meta := &sheetMeta{
				Merges: []*sheetMetaMerge{
					{
						ColumnCount:      2,
						RowCount:         2,
						StartColumnIndex: 0,
						StartRowIndex:    0,
					},
				},
			}
			p := gomonkey.ApplyPrivateMethod(mSheet, "getMeta", func(_ *Sheet) (*sheetMeta, error) {
				return meta, nil
			})
			defer p.Reset()
			rows, err := content.ToRowsParseMerged()
			So(err, ShouldBeNil)
			So(rows, ShouldResemble, []SheetRow{
				{NewSheetCell("ace"), NewSheetCell("ace")},
				{NewSheetCell("ace"), NewSheetCell("ace")},
			})
		})
	})
}

func TestWriteAndReadData(t *testing.T) {
	sheet := getSheetClient().AddSheet("test write and read", 0)
	assert.NoError(t, sheet.Err)
	sheet.NewRangeFull("F2", "F2").SetDropdown([]string{"Monday", "Tuesday", "Wednesday"}, true, nil)
	sheet.WriteRows([][]interface{}{
		{"string", "no text link", "link", "email", "user", "formula"},
		{"test string", "https://z.cn", SheetCellTypeLink("amazon", "https://z.cn"),
			SheetCellTypeMentionEmail("", false, true), SheetCellTypeFormula("=A1"),
			SheetCellTypeDropdown([]interface{}{"Monday", "Tuesday"}),
		},
	})
	assert.NoError(t, sheet.Err)
	_, err := sheet.ReadRows()
	assert.NoError(t, err)
	/*
		rows = rows.To TrimBlankTail()
		for _, row := range rows {
			for _, cell := range row {
				t.Log(cell.Value())
			}
		}
	*/
	rows, err := sheet.ReadRows()
	assert.NoError(t, err)
	rows = sheet.TrimBlankTail(rows)
	for _, row := range rows {
		for _, cell := range row {
			t.Log(cell.Value())
		}
	}
}

func TestSheet_moveRowOrColumn(t *testing.T) {
	getSheet()
	sheet.WriteRows([][]interface{}{
		{"a", "b", "c"},
		{"d", "e", "f"},
		{"g", "h", "k"},
	})
	assert.NoError(t, sheet.Err)
	sheet.MoveRows(1, 1, 3)
	assert.NoError(t, sheet.Err)
	sheet.MoveColumns(1, 1, 5)
	assert.NoError(t, sheet.Err)
}

func TestSheet_update(t *testing.T) {
	getSheet()
	assert.NoError(t, sheet.Err)
	meta, err := sheet.getMeta()
	assert.NoError(t, err)
	assert.Equal(t, meta.Title, "test update")
	assert.Equal(t, meta.Index, 0)

	t.Run("title", func(t *testing.T) {
		sheet.UpdateTitle("test update 1")
		meta, _ := sheet.getMeta()
		assert.Equal(t, meta.Title, "test update 1")
	})
	t.Run("index", func(t *testing.T) {
		sheet.UpdateIndex(2)
		meta, _ := sheet.getMeta()
		assert.Equal(t, meta.Index, 2)
	})
	t.Run("hidden", func(t *testing.T) {
		sheet = sheet.Hidden(true)
		assert.NoError(t, sheet.Err)
		sheet = sheet.Hidden(false)
	})
	t.Run("frozen", func(t *testing.T) {
		sheet.FrozenColumn(3)
		sheet.FrozenRow(3)
		assert.NoError(t, sheet.Err)
	})
	t.Run("protect", func(t *testing.T) {
		sheet.Protect("locked", nil)
		assert.NoError(t, sheet.Err)
	})
}

func getSheet() *Sheet {
	if sheet == nil {
		sheet = getSheetClient().AddSheet("test update", 0)
	}
	return sheet
}
