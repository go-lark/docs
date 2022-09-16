package docs

import (
	"fmt"
	"testing"

	"github.com/hilaily/kit/stringx"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestSpreadSheetOrigin_MetaInfo(t *testing.T) {
	Convey("TestSpreadSheetOrigin_MetaInfo", t, func() {
		b, res, err := getOrigin().MetaInfo()
		assert.NoError(t, err)
		So(res, ShouldNotBeNil)
		assert.NotEmpty(t, res.SpreadsheetToken)
		assert.Equal(t, res.SpreadsheetToken, testSpreadSheetToken)
		t.Log(string(b))
	})
}

func TestSpreadSheetOrigin_Properties(t *testing.T) {
	_, res, err := getOrigin().MetaInfo()
	assert.NoError(t, err)
	assert.NotNil(t, res)
	title := res.Properties.Title
	t.Log("old title: ", title)
	newTitle := "new title " + stringx.GenRankStr(10)
	_, err = getOrigin().Properties(&SpreadSheetProperties{Title: &newTitle})
	assert.NoError(t, err)
	_, res, _ = getOrigin().MetaInfo()
	assert.Equal(t, newTitle, res.Properties.Title)
	assert.NotEqual(t, title, res.Properties.Title)
}

func TestSpreadSheetOrigin_SheetBatchUpdate(t *testing.T) {
	title := "t" + stringx.GenRankStr(5)
	res, _, err := getOrigin().SheetBatchUpdate(map[ModifySheetType]interface{}{
		ModifySheetAdd: map[string]interface{}{
			"properties": map[string]interface{}{
				"title": title,
			},
		},
	})
	assert.NoError(t, err)
	t.Log(res)
}

/*
func TestSpreadSheetOrigin_ValuesPrepend(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	r := NewRangeCol(id, "A1", "B2")
	t.Log("r: ", r)
	data := [][]interface{}{{"string", 1}, {"haha", "hehe"}}
	b, res2, err := getOrigin().ValuesPrepend(r, data)
	assert.NoError(t, err)
	assert.NotZero(t, len(b))
	t.Log(res2)
	b, res2, err = getOrigin().ValuesAppend(NewRangeFull(id, "A3", "B4"), data, InseartDataOptionOverwrite)
	assert.NoError(t, err)
	assert.NotZero(t, len(b))
	t.Log(res2)
}
*/

func TestSpreadSheetOrigin_InsertiDimensionRange(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	err = getOrigin().InsertDimensionRange(id, 3, 5, MajorDimensionRows, InseartDataOptionOverwrite)
	assert.NoError(t, err)
}

func TestSpreadSheetOrigin_DimensionRange(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	_, _, err = getOrigin().DimensionRangePost(id, MajorDimensionRows, 10)
	assert.NoError(t, err)
}

func TestSpreadSheetOrigin_DimensionRangePut(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	_, _, err = getOrigin().DimensionRangePut(id, MajorDimensionRows, 1, 10, true, 10)
	assert.NoError(t, err)
}

func TestSpreadSheetOrigin_DimensionRangeDelete(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	_, _, err = getOrigin().DimensionRangeDelete(id, MajorDimensionRows, 1, 1)
	assert.NoError(t, err)
}

/*
func TestSpreadSheetOrigin_WriteValuesByRange(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	b, _, err := getOrigin().WriteValuesByRange(NewRangeFull(id, "A1", "C3"), [][]interface{}{{"1", "2", "3"}, {4, 5, 6}})
	assert.NoError(t, err)
	assert.NotEmpty(t, string(b))
	t.Log(string(b))
}

func TestSpreadSheetOrigin_ReadValuesByRange(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	_, res, err := getOrigin().ReadValuesByRange(NewRangeFull(id, "A1", "B2"), "", "")
	assert.NoError(t, err)
	assert.NotZero(t, len(res.ValueRange.Values))
	t.Log(res.ValueRange.Values)
}

func TestSpreadSheetOrigin_WriteValuesByRangeMulti(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	b, _, err := getOrigin().WriteValuesByRangeMulti(
		[]*WriteValuesByRangeMultiArgs{
			{NewRangeFull(id, "A1", "C3"), [][]interface{}{{"1", "2", "3"}, {4, 5, 6}}},
			{NewRangeFull(id, "A4", "C6"), [][]interface{}{{"1", "2", "3"}, {4, 5, 6}}},
		},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, string(b))
	t.Log(string(b))
}

func TestSpreadSheetOrigin_ReadValuesByRangeMulti(t *testing.T) {
	id, err := getSheetID()
	assert.NoError(t, err)
	_, res, err := getOrigin().ReadValuesByRangeMulti([]Range{
		NewRangeFull(id, "A1", "C6"),
		NewRangeFull(id, "A4", "C6"),
	},
		"", "",
	)
	assert.NoError(t, err)
	assert.NotZero(t, len(res.ValueRanges))
	t.Log(res.ValueRanges)
}
*/

func getOrigin() *SpreadSheetOrigin {
	return &SpreadSheetOrigin{
		baseClient: getClientNew(),
		token:      testSpreadSheetToken,
	}
}

func getSheetID() (string, error) {
	_, res, err := getOrigin().MetaInfo()
	if err != nil {
		return "", err
	}
	if len(res.Sheets) == 0 {
		return "", fmt.Errorf("sheets is empty")
	}
	return res.Sheets[0].SheetID, nil
}
