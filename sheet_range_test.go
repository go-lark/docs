package docs

import (
	"testing"

	"github.com/hilaily/kit/dev"
	"github.com/stretchr/testify/assert"
)

func TestPrepend(t *testing.T) {
	c := getSheetClient().GetSheetByIndex(0)
	r := c.GetRange("A1", "B2").Prepend(
		[][]interface{}{
			{"a", "b"},
			{1, 3},
		},
	)
	assert.NoError(t, r.Err)
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
