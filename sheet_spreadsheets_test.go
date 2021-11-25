package docs

import (
	"fmt"
	"testing"
	"time"

	"github.com/hilaily/kit/stringx"
	"github.com/stretchr/testify/assert"
)

func TestSpreadSheets_ChangeOwner(t *testing.T) {
	ss := getClientNew().RootFolder().CreateSpreadSheet("test create sheet"+time.Now().String()).ChangeOwner(NewMemberWithEmail(testUserEmail), false, false)
	assert.NoError(t, ss.Err)
	t.Log(baseDomain + "/sheets/" + ss.token)
	sheet := ss.SheetIndex(0).WriteRows([]string{"name", "age"}, [][]interface{}{
		{"Ace", 1}, {"Bob", 2},
	})
	assert.NoError(t, sheet.Err)
}

func TestSpreadSheets_AddMember(t *testing.T) {
	ss := getClientNew().RootFolder().CreateSpreadSheet("test create sheet"+time.Now().String()).
		Share(PermEdit, false, NewMemberWithEmail(testUserEmail))
	assert.NoError(t, ss.Err)
	t.Log(ss.token)
}

func TestSpreadSheets_Meta(t *testing.T) {
	sheet := getClientNew().OpenSpreadSheets(testSpreadSheetToken)
	meta, err := sheet.GetMeta()
	assert.NoError(t, err)
	assert.NotZero(t, len(meta.Sheets))
}

func TestSpreadSheets_Content(t *testing.T) {
	c := getSheetClient().SheetID("f6d5a1")
	res, err := c.GetContentByRange("A1", "A1")
	assert.Nil(t, err)
	assert.NotZero(t, len(res.ToRows()))
}

func TestSpreadSheets_V2(t *testing.T) {
	c := getSheetClient().SheetID("2BGf04")
	res, err := c.GetContentByRangeV2("A1", "A1", SheetRenderFormula, "")
	assert.Nil(t, err)
	assert.NotZero(t, len(res.ToRows()))
	t.Log(res)
}

func TestSpreadSheets_GetContent(t *testing.T) {
	c := getSheetClient().SheetID("2BGf04")
	rows, err := c.GetRows(false)
	assert.NoError(t, err)
	assert.NotZero(t, len(rows))
	t.Log("count: ", len(rows))
}

func TestSpreadSheets_WriteRows(t *testing.T) {
	sheet := getSheetClient().SheetName("Sheet1")
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
	sheet := getSheetClient().SheetName("Sheet1")
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

func TestNow(t *testing.T) {
	now := time.Now().Unix()
	t.Log(now)
}

func TestSpreadSheets_ColnameAdd(t *testing.T) {
	data := []struct {
		start string
		add   int
		val   string
	}{
		{"A", 1, "B"},
		{"C", 4, "G"},
		{"Z", 1, "AA"},
		{"Z", 2, "AB"},
		{"ABC", 27, "ACD"},
		{"A", 100, "CW"},
	}
	for _, v := range data {
		res := colnameAdd(v.start, v.add)
		assert.Equal(t, v.val, res, fmt.Sprintf("start: %s", v.start))
	}
}

func TestSpreadSheets_ColnameSplit(t *testing.T) {
	data := []struct {
		val    string
		expCol string
		expRow int
	}{
		{"A1", "A", 1},
		{"A12", "A", 12},
		{"BC43", "BC", 43},
		{"A401", "A", 401},
	}
	for _, v := range data {
		col, row := cellnameSplit(v.val)
		assert.Equal(t, v.expCol, col, fmt.Sprintf("cell: %s", v.val))
		assert.Equal(t, v.expRow, row, fmt.Sprintf("cell: %s", v.val))
	}
}

func TestSpreadSheets_UpdateTitle(t *testing.T) {
	ss := getSheetClient().UpdateTitle("update title")
	assert.Nil(t, ss.Err)
}

func TestSpreadSheets_AddSheet(t *testing.T) {
	title := "t" + stringx.GenRankStr(5)
	sheet := getSheetClient().AddSheet(title, 0)
	assert.Nil(t, sheet.Err)
	assert.NotEmpty(t, sheet.id)
	newSheetID := sheet.GetID()
	sheet = getSheetClient().CopySheet("4gJAV3", "t"+stringx.GenRankStr(5))
	assert.Nil(t, sheet.Err)
	assert.NotEmpty(t, sheet.id)
	ss := getSheetClient().DeleteSheet(newSheetID)
	assert.Nil(t, ss.Err)
}

func getSheetClient() *SpreadSheets {
	sheet := getClientNew().OpenSpreadSheets(testSpreadSheetToken)
	return sheet
}

type row struct {
	Title string
	Task  string
	Date  time.Time
	Mark  string
}
