package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDropdown(t *testing.T) {
	sheet := getSheetClient().AddSheet("test set dropdown", 0)
	assert.NoError(t, sheet.Err)
	err := sheet.NewRangeFull("A1", "A10").SetDropdown([]string{"1", "2", "3"}, true, []string{"#1FB6C1", "#F006C2", "#CB1AC3"})
	assert.NoError(t, err)
	ss := sheet.ssClient.DeleteSheet(sheet.GetID())
	assert.NoError(t, ss.Err)
}

/*
func TestGetRows(t *testing.T) {
	c := getSheetClient().SheetID("2BGf04")
	c.Scan(nil)
}
*/
