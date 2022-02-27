package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sheet *Sheet
)

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
