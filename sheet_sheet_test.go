package docs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sheet *Sheet
)

func TestSpreadSheetsBind(t *testing.T) {
	c := getSpreadSheet().GetSheetByID("2BGf04")
	rows, err := c.GetRows(true)
	assert.NoError(t, err)
	assert.NotZero(t, len(rows))
	t.Log("count: ", len(rows))
	for i, v := range rows {
		for j, vv := range v {
			t.Logf("%d, %d, %v", i, j, vv.ToString())
		}
	}
}

func TestColCul(t *testing.T) {
	data := map[string]int{
		"A":  1,
		"Z":  26,
		"AA": 27,
	}
	for k, v := range data {
		r := colName2Num(k)
		assert.Equal(t, v, r, k)
	}
}

func TestSpreadSheets_Meta(t *testing.T) {
	sheet := getClientNew().OpenSpreadSheet(testSpreadSheetToken)
	meta, err := sheet.GetMeta()
	assert.NoError(t, err)
	assert.NotZero(t, len(meta.Sheets))
	t.Log(meta.Sheets[0].Merges[0])
}

func TestSpreadSheets_Content(t *testing.T) {
	c := getSpreadSheet().GetSheetByID("f6d5a1")
	res, err := c.GetContentByRange("A1", "A1")
	assert.Nil(t, err)
	assert.NotZero(t, len(res.ToRows()))
}

func TestSpreadSheets_V2(t *testing.T) {
	c := getSpreadSheet().GetSheetByID("2BGf04")
	res, err := c.GetContentByRangeV2("A1", "A1", SheetRenderFormula, "")
	assert.Nil(t, err)
	assert.NotZero(t, len(res.ToRows()))
	t.Log(res)
}

func TestSpreadSheets_GetContent(t *testing.T) {
	c := getSpreadSheet().GetSheetByID("2BGf04")
	rows, err := c.GetRows(true)
	assert.NoError(t, err)
	assert.NotZero(t, len(rows))
	t.Log("count: ", len(rows))
	for i, v := range rows {
		for j, vv := range v {
			t.Logf("%d, %d, %v", i, j, vv.ToString())
		}
	}
}

func TestSpreadSheets_WriteRows(t *testing.T) {
	sheet := getSpreadSheet().GetSheetByName("Sheet1")
	assert.NoError(t, sheet.Err)
	title := []string{"first col", "second col", "third col"}
	data := [][]interface{}{
		{"1", "2", "3"},
		{1, 2, 3},
		{4, nil, 6},
		{7, "", 9},
	}
	sheet = sheet.WriteRows(title, nil)
	assert.Nil(t, sheet.Err)
	sheet = sheet.WriteRows(title, data)
	assert.Nil(t, sheet.Err)
}

func TestSpreadSheets_WriteALotRows(t *testing.T) {
	sheet := getSpreadSheet().GetSheetByName("Sheet1")
	assert.NoError(t, sheet.Err)
	data := [][]interface{}{}
	colCount := 10
	rowCount := 10
	for i := 0; i < rowCount; i++ {
		d := make([]interface{}, 0, colCount)
		for j := 0; j < colCount; j++ {
			if j == 0 {
				d = append(d, i)
			} else {
				d = append(d, j)
			}
		}
		data = append(data, d)
	}
	sheet = sheet.WriteRows(nil, data)
	assert.Nil(t, sheet.Err)
}

func TestASCII(t *testing.T) {
	fmt.Printf("%d\n", 'A')
	fmt.Printf("%d\n", 'Z')
	fmt.Println(int('A'))
	fmt.Println(string(byte(66)))
}

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
		sheet = getSpreadSheet().CreateSheet("test update", 0)
	}
	return sheet
}
