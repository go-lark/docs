/*
	SpreadSheetOrigin provide the API of the feihsu/lark open API document.
	SpreadSheet provide the combinatorial API of SpreadSheetOrigin.
*/

package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hilaily/kit/netx"
	"github.com/hilaily/kit/stringx"
)

func newSpreadSheetOrigin(client *Client, token string) *SpreadSheetOrigin {
	return &SpreadSheetOrigin{
		baseClient: client,
		token:      token,
	}
}

// SpreadSheetOrigin represent for origin open API.
// 这个客户端的方法使用了原始的开放平台 API，没有做任何封装。
type SpreadSheetOrigin struct {
	baseClient *Client
	token      string
}

// MetaInfo
// api: /open-apis/sheet/v2/spreadsheets/:spreadsheetToken/metainfo
// reference: https://open.feishu.cn/document/ukTMukTMukTM/uETMzUjLxEzM14SMxMTN
func (so *SpreadSheetOrigin) MetaInfo() (b []byte, result *MetaInfoResp, err error) {
	url := so.baseClient.urlJoin("/open-apis/sheet/v2/spreadsheets/", so.token, "metainfo")
	_req, _ := http.NewRequest(http.MethodGet, url, nil)
	b, err = so.baseClient.CommonReq(_req, &result)
	return
}

// Properties
// reference: https://open.feishu.cn/document/ukTMukTMukTM/ucTMzUjL3EzM14yNxMTN
func (so *SpreadSheetOrigin) Properties(prop *SpreadSheetProperties) (b []byte, err error) {
	en, _ := json.Marshal(map[string]*SpreadSheetProperties{
		"properties": prop,
	})
	req, _ := http.NewRequest(
		http.MethodPut,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "/properties"),
		bytes.NewReader(en),
	)
	b, err = so.baseClient.CommonReq(req, nil)
	return
}

// SheetBatchUpdate for update properties of a sheet
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uQDO2UjL0gjN14CN4YTN
// 更新 sheet 的属性
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/ugjMzUjL4IzM14COyMTN
func (so *SpreadSheetOrigin) SheetBatchUpdate(args map[ModifySheetType]interface{}) (b []byte, resp *SheetBatchUpdateResp, err error) {
	argsArr := make([]interface{}, 0, 1)
	for k := range args {
		argsArr = append(argsArr, map[ModifySheetType]interface{}{
			k: args[k],
		})
	}

	en, _ := json.Marshal(map[string][]interface{}{
		"requests": argsArr,
	})
	resp = &SheetBatchUpdateResp{}
	_req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "open-apis/sheet/v2/spreadsheets/", so.token, "/sheets_batch_update"),
		bytes.NewReader(en))
	b, err = so.baseClient.CommonReq(_req, resp)
	return
}

// Import represent import excel from local
// reference https://open.larksuite.com/document/ukTMukTMukTM/uATO2YjLwkjN24CM5YjN
func (so *SpreadSheetOrigin) Import(filepath string, filename string, folderToken string) ([]byte, string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, "", err
	}
	en, _ := json.Marshal(
		map[string]interface{}{
			"file":        data,
			"name":        filename,
			"folderToken": folderToken,
		},
	)
	req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheets/v2/import"),
		bytes.NewReader(en),
	)
	resp := &struct {
		Ticket string `json:"ticket"`
	}{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp.Ticket, err
}

// ValuesPrepend
// reference https://open.feishu.cn/document/ukTMukTMukTM/uIjMzUjLyIzM14iMyMTN
func (so *SpreadSheetOrigin) ValuesPrepend(_range SheetRange, data [][]interface{}) ([]byte, *ValuesPrependResp, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"valueRange": map[string]interface{}{
			"range":  _range,
			"values": data,
		},
	})
	req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "values_prepend"),
		bytes.NewReader(en),
	)
	resp := &ValuesPrependResp{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp, err
}

// ValuesAppend
// reference https://open.feishu.cn/document/ukTMukTMukTM/uMjMzUjLzIzM14yMyMTN
func (so *SpreadSheetOrigin) ValuesAppend(_range Range, data [][]interface{}, inseartDataOption InseartDataOptionType) ([]byte, *ValuesPrependResp, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"valueRange": map[string]interface{}{
			"range":  _range,
			"values": data,
		},
	})
	req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "/values_prepend?insertDataOption="+inseartDataOption),
		bytes.NewReader(en),
	)
	resp := &ValuesPrependResp{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp, err
}

func (so *SpreadSheetOrigin) InsertDimensionRange(sheetID string, startIndex, endIndex int, majorDimension MajorDimensionType, inheritStyle InheritStyleType) error {
	en, _ := json.Marshal(map[string]interface{}{
		"dimension": map[string]interface{}{
			"sheetId":        sheetID,
			"majorDimension": majorDimension,
			"startIndex":     startIndex,
			"endIndex":       endIndex,
		},
		"inheritStyle": inheritStyle,
	})
	req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "/insert_dimension_range"),
		bytes.NewReader(en),
	)
	_, err := so.baseClient.CommonReq(req, nil)
	return err

}

func (so *SpreadSheetOrigin) DimensionRangePost(sheetID string, majorDimension MajorDimensionType, length int) ([]byte, *DimensionRangeResp, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"dimension": map[string]interface{}{
			"sheetId":        sheetID,
			"majorDimension": majorDimension,
			"length":         length,
		},
	})
	req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "/dimension_range"),
		bytes.NewReader(en),
	)
	resp := &DimensionRangeResp{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp, err
}

// DimensionRangePut for update rows or columns
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uETO2UjLxkjN14SM5YTN
// 更新行列
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uYjMzUjL2IzM14iNyMTN
func (so *SpreadSheetOrigin) DimensionRangePut(sheetID string, majorDimension MajorDimensionType, startIndex, endIndex int, visible bool, fixedSize int) ([]byte, *DimensionRangePutResp, error) {
	en, _ := json.Marshal(map[string]map[string]interface{}{
		"dimension": {
			"sheetId":        sheetID,
			"majorDimension": majorDimension,
			"startIndex":     startIndex,
			"endIndex":       endIndex,
		},
		"dimensionProperties": {
			"visible":   visible,
			"fixedSize": fixedSize,
		},
	})
	req, _ := http.NewRequest(
		http.MethodPut,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "/dimension_range"),
		bytes.NewReader(en),
	)
	b, err := so.baseClient.CommonReq(req, nil)
	return b, nil, err
}

// DimensionRangeDelete for delete rows or columns
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uATO2UjLwkjN14CM5YTN
// 删除行列
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/ucjMzUjL3IzM14yNyMTN
func (so *SpreadSheetOrigin) DimensionRangeDelete(sheetID string, majorDimension MajorDimensionType, startIndex, endIndex int) ([]byte, *DimensionRangeDeleteResp, error) {
	en, _ := json.Marshal(map[string]map[string]interface{}{
		"dimension": {
			"sheetId":        sheetID,
			"majorDimension": majorDimension,
			"startIndex":     startIndex,
			"endIndex":       endIndex,
		},
	})
	req, _ := http.NewRequest(
		http.MethodDelete,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheet/v2/spreadsheets/", so.token, "/dimension_range"),
		bytes.NewReader(en),
	)
	resp := &DimensionRangeDeleteResp{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp, err
}

// ReadValuesByRange read data by range
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/ukzN2UjL5cjN14SO3YTN
// 读取单个范围的数据
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/ugTMzUjL4EzM14COxMTN
func (so *SpreadSheetOrigin) ReadValuesByRange(_range Range, valueRender string, dateTimeRender string) ([]byte, *GetValuesByRangeResp, error) {
	u := so.baseClient.urlJoin("open-apis/sheet/v2/spreadsheets/", so.token, "/values/", _range)
	p := map[string]string{}
	if valueRender != "" {
		p["valueRenderOption"] = valueRender
	}
	if dateTimeRender != "" {
		p["dateTimeRenderOption"] = dateTimeRender
	}
	u, _ = netx.URLQueryParams(u, p)
	_req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp := &GetValuesByRangeResp{}
	b, err := so.baseClient.CommonReq(_req, resp)
	return b, resp, err
}

// WriteValuesByRange for writing data to a single range
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uMDO2UjLzgjN14yM4YTN
// 向单个范围写入数据
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uAjMzUjLwIzM14CMyMTN
func (so *SpreadSheetOrigin) WriteValuesByRange(_range Range, data [][]interface{}) ([]byte, *WriteValuesByRangeResp, error) {
	en, _ := json.Marshal(map[string]map[string]interface{}{
		"valueRange": {
			"range":  _range,
			"values": data,
		},
	})
	req, _ := http.NewRequest(
		http.MethodPut,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheets/v2/spreadsheets/", so.token, "/values"),
		bytes.NewReader(en),
	)
	resp := &WriteValuesByRangeResp{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp, err
}

// ReadValuesByRangeMulti for reading multiple ranges
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uADO2UjLwgjN14CM4YTN
// 读取多个范围的数据
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/ukTMzUjL5EzM14SOxMTN
func (so *SpreadSheetOrigin) ReadValuesByRangeMulti(ranges []Range, valueRender string, dateTimeRender string) ([]byte, *GetValuesByRangeMultiResp, error) {
	u := so.baseClient.urlJoin("open-apis/sheet/v2/spreadsheets/", so.token, "/values_batch_get/")
	p := map[string]string{
		"ranges": strings.Join(ranges, ","),
	}
	if valueRender != "" {
		p["valueRenderOption"] = valueRender
	}
	if dateTimeRender != "" {
		p["dateTimeRenderOption"] = dateTimeRender
	}
	u, _ = netx.URLQueryParams(u, p)
	_req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp := &GetValuesByRangeMultiResp{}
	b, err := so.baseClient.CommonReq(_req, resp)
	return b, resp, err
}

// WriteValuesByRangeMulti for writing data to multiple ranges
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uIDO2UjLygjN14iM4YTN
// 向多个范围写入数据
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uEjMzUjLxIzM14SMyMTN
func (so *SpreadSheetOrigin) WriteValuesByRangeMulti(rangeDatas []*WriteValuesByRangeMultiArgs) ([]byte, *WriteValuesByRangeMultiResp, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"valueRanges": rangeDatas,
	})
	req, _ := http.NewRequest(
		http.MethodPost,
		stringx.URLJoin(so.baseClient.domain, "/open-apis/sheets/v2/spreadsheets/", so.token, "values_batch_update"),
		bytes.NewReader(en),
	)
	resp := &WriteValuesByRangeMultiResp{}
	b, err := so.baseClient.CommonReq(req, resp)
	return b, resp, err
}

// Style ...
func (so *SpreadSheetOrigin) Style(_range Range, style *SheetCellStyle) (*SheetStyleResp, error) {
	u := so.baseClient.urlJoin("/open-apis/sheets/v2/spreadsheets/" + so.token + "/style")
	en, _ := json.Marshal(map[string]interface{}{
		"appendStyle": map[string]interface{}{
			"range": _range,
			"style": style,
		},
	})
	req, _ := http.NewRequest(http.MethodPut,
		u,
		bytes.NewReader(en),
	)
	r := &SheetStyleResp{}
	_, err := so.baseClient.CommonReq(req, r)
	if err != nil {
		return nil, fmt.Errorf("set style fail %w", err)
	}
	return r, nil
}

type SpreadSheetProperties struct {
	Title *string `json:"title,omitempty"`
}

type ValuesPrependResp struct {
	Revision         int                  `json:"revision"`
	SpreadsheetToken string               `json:"spreadsheetToken"`
	TableRange       string               `json:"tableRange"`
	Updates          ValuesPrependUpdates `json:"updates"`
}

type ValuesPrependUpdates struct {
	SpreadsheetToken string `json:"spreadsheetToken"`
	UpdatedRange     string `json:"updatedRange"`
	UpdatedRows      int    `json:"updatedRows"`
	UpdatedColumns   int    `json:"updatedColumns"`
	UpdatedCells     int    `json:"updatedCells"`
	Revision         int    `json:"revision"`
}

type DimensionRangeResp struct {
	AddCount       int    `json:"addCount"`
	MajorDimension string `json:"majorDimension"`
}

type (
	SheetBatchUpdateResp struct {
		Replies []SheetBatchUpdateRespReplies `json:"replies"`
	}

	SheetBatchUpdateRespReplies struct {
		AddSheet    SheetBatchUpdateRespAddSheet    `json:"addSheet"`
		CopySheet   SheetBatchUpdateRespCopySheet   `json:"copySheet"`
		UpdateSheet SheetBatchUpdateRespUpdateSheet `json:"updateSheet"`
		DeleteSheet SheetBatchUpdateRespDeleteSheet `json:"deleteSheet"`
	}

	SheetBatchUpdateRespProperties struct {
		SheetID string `json:"sheetId"`
		Title   string `json:"title"`
		Index   int    `json:"index"`
	}

	SheetBatchUpdateRespAddSheet struct {
		Properties SheetBatchUpdateRespProperties `json:"properties"`
	}
	SheetBatchUpdateRespCopySheet struct {
		Properties SheetBatchUpdateRespProperties `json:"properties"`
	}

	SheetBatchUpdateRespDeleteSheet struct {
		Result  bool   `json:"result"`
		SheetID string `json:"sheetId"`
	}

	SheetBatchUpdateRespUpdateSheet struct {
		UpdateSheet SheetBatchUpdateRespUpdateSheetUpdateSheet `json:"updateSheet"`
	}

	SheetBatchUpdateRespUpdateSheetProtect struct {
		Lock      string  `json:"lock"`
		SheetName string  `json:"sheetName"`
		PermID    string  `json:"permId"`
		UserIds   []int64 `json:"userIds"`
	}

	SheetBatchUpdateRespUpdateSheetProperties struct {
		SheetID        string                                 `json:"sheetId"`
		Title          string                                 `json:"title"`
		Index          int                                    `json:"index"`
		Hidden         bool                                   `json:"hidden"`
		FrozenColCount int                                    `json:"frozenColCount"`
		FrozenRowCount int                                    `json:"frozenRowCount"`
		Protect        SheetBatchUpdateRespUpdateSheetProtect `json:"protect"`
	}

	SheetBatchUpdateRespUpdateSheetUpdateSheet struct {
		Properties SheetBatchUpdateRespUpdateSheetProperties `json:"properties"`
	}
)

type (
	DimensionRangePutResp = interface{} // there is no document, for future use.
)

type (
	DimensionRangeDeleteResp struct {
		DelCount       int                `json:"delCount"`
		MajorDimension MajorDimensionType `json:"majorDimension"`
	}
)

type (
	GetValuesByRangeResp struct {
		Revision         int                            `json:"revision"`
		SpreadsheetToken string                         `json:"spreadsheetToken"`
		ValueRange       GetValuesByRangeRespValueRange `json:"valueRange"`
	}

	GetValuesByRangeRespValueRange struct {
		MajorDimension string          `json:"majorDimension"`
		Range          string          `json:"range"`
		Revision       int             `json:"revision"`
		Values         [][]interface{} `json:"values"`
	}
)
type (
	GetValuesByRangeMultiResp struct {
		Revision         int                                    `json:"revision"`
		SpreadsheetToken string                                 `json:"spreadsheetToken"`
		TotalCells       int                                    `json:"totalCells"`
		ValueRanges      []GetValuesByRangeMultiRespValueRanges `json:"valueRanges"`
	}

	GetValuesByRangeMultiRespValueRanges struct {
		MajorDimension string          `json:"majorDimension"`
		Range          string          `json:"range"`
		Revision       int             `json:"revision"`
		Values         [][]interface{} `json:"values"`
	}
)

type (
	WriteValuesByRangeResp struct {
		Revision         int    `json:"revision"`
		SpreadsheetToken string `json:"spreadsheetToken"`
		UpdatedCells     int    `json:"updatedCells"`
		UpdatedColumns   int    `json:"updatedColumns"`
		UpdatedRange     string `json:"updatedRange"`
		UpdatedRows      int    `json:"updatedRows"`
	}
)
type (
	WriteValuesByRangeMultiArgs struct {
		Range  Range           `json:"range"`
		Values [][]interface{} `json:"values"`
	}

	WriteValuesByRangeMultiResp struct {
		Responses        []WriteValuesByRangeMultiRespResponses `json:"responses"`
		Revision         int                                    `json:"revision"`
		SpreadsheetToken string                                 `json:"spreadsheetToken"`
	}

	WriteValuesByRangeMultiRespResponses struct {
		SpreadsheetToken string `json:"spreadsheetToken"`
		UpdatedCells     int    `json:"updatedCells"`
		UpdatedColumns   int    `json:"updatedColumns"`
		UpdatedRange     string `json:"updatedRange"`
		UpdatedRows      int    `json:"updatedRows"`
	}
)

type (
	SheetStyleResp struct {
		Updates SheetStyleRespUpdate `json:"updates"`
	}
	SheetStyleRespUpdate struct {
		SpreadsheetToken string `json:"spreadsheetToken"`
		UpdatedRange     string `json:"updatedRange"`
		UpdatedRows      int    `json:"updatedRows"`
		UpdatedColumns   int    `json:"updatedColumns"`
		UpdatedCells     int    `json:"updatedCells"`
		Revision         int    `json:"revision"`
	}
)
