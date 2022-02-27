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
	sheet := ss.GetSheetByIndex(0).WriteRows(
		[]string{"name", "age"},
		[][]interface{}{
			{"Ace", 1},
			{"Bob", 2},
			{"", ""},
		},
	)
	assert.NoError(t, sheet.Err)
}

func TestSpreadSheets_AddMember(t *testing.T) {
	spreadSheet := getClient().RootFolder().CreateSpreadSheet("create sheet"+time.Now().String()).ChangeOwner(NewMemberWithEmail(testUserEmail), false, false)
	assert.NoError(t, spreadSheet.Err)
	u := NewMemberWithUserID(testUserID)
	spreadSheet = spreadSheet.Share(PermEdit, false, u)
	assert.NoError(t, spreadSheet.Err)
	t.Log(spreadSheet.GetToken())
	spreadSheet = spreadSheet.UnShare(u)
	assert.NoError(t, spreadSheet.Err)
	//err := spreadSheet.DeleteSelf()
	//assert.NoError(t, err)
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

func TestGetSheet(t *testing.T) {
	spreadSheet := getSpreadSheet()
	sheet := spreadSheet.GetSheetByIndex(0)
	assert.NoError(t, sheet.Err)
	assert.NotEmpty(t, sheet.GetID())
	sheet = spreadSheet.GetSheetByID(sheet.GetID())
	assert.NoError(t, sheet.Err)
	assert.NotEmpty(t, sheet.GetID())
	sheet = spreadSheet.GetSheetByName(sheet.GetName())
	assert.NoError(t, sheet.Err)
	assert.NotEmpty(t, sheet.GetID())
}

func TestSpreadSheets_UpdateTitle(t *testing.T) {
	meta, err := getSpreadSheet().GetMeta()
	assert.NoError(t, err)
	oldTtitle := meta.Properties.Title
	newTitle := "update title"
	ss := getSpreadSheet().UpdateTitle(newTitle)
	assert.NoError(t, ss.Err)
	meta, err = getSpreadSheet().GetMeta()
	assert.NoError(t, err)
	assert.Equal(t, newTitle, meta.Properties.Title)
	ss = getSpreadSheet().UpdateTitle(oldTtitle)
	assert.NoError(t, ss.Err)
}

func TestSpreadSheets_AddSheet(t *testing.T) {
	title := "t" + stringx.GenRankStr(5)
	sheet := getSpreadSheet().CreateSheet(title, 0)
	assert.NoError(t, sheet.Err)
	assert.NotEmpty(t, sheet.GetID())
	newSheetID := sheet.GetID()
	newTitle := "copy a new sheet"
	sheet = getSpreadSheet().CopySheet(sheet.GetID(), newTitle)
	assert.NoError(t, sheet.Err)
	assert.NotEmpty(t, sheet.GetID())
	ss := getSpreadSheet().DeleteSheet(newSheetID)
	assert.NoError(t, ss.Err)
	ss = getSpreadSheet().DeleteSheet(sheet.GetID())
	assert.NoError(t, ss.Err)
}

func TestCreateSpreadSheet(t *testing.T) {
	spreadSheet := getClient().RootFolder().CreateSpreadSheet("a test sheet")
	assert.NoError(t, spreadSheet.Err)
	t.Log(spreadSheet.GetToken())
}

func getSpreadSheet() *SpreadSheet {
	sheet := getClient().OpenSpreadSheet(testSpreadSheetToken)
	return sheet
}

type row struct {
	Title string
	Task  string
	Date  time.Time
	Mark  string
}
