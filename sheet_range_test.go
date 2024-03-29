package docs

import (
	"testing"

	"github.com/hilaily/kit/dev"
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

func TestScan(t *testing.T) {
	type A struct {
		Name string
		Age  int
	}
	rows := []SheetRow{
		[]*SheetCell{NewSheetCell("ace"), NewSheetCell(2)},
		[]*SheetCell{NewSheetCell("bob"), NewSheetCell(3)},
	}
	r := &SheetRange{}
	res := []*A{}
	err := r.scan(rows, &res)
	assert.NoError(t, err)
	dev.PJSON(res)
}

/*
func TestGetRows(t *testing.T) {
	c := getSheetClient().SheetID("2BGf04")
	c.Scan(nil)
}
*/
