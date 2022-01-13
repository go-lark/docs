package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sheet *Sheet
)

func TestSheet_moveRowOrColumn(t *testing.T) {
	getSheet()
	sheet.WriteRows([]string{"a", "b", "c"}, [][]interface{}{
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
