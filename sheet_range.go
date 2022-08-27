package docs

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Range reference https://open.feishu.cn/document/ukTMukTMukTM/uczNzUjL3czM14yN3MTN#bae19f77
type SheetRange struct {
	Err         error
	leftTop     *cellName
	rightBottom *cellName
	sheet       *Sheet
}

/*
// Prepend insert data before range block.
func (s *SheetRange) Prepend() {}

// Append insert data after range block.
func (s *SheetRange) Append() {}

// Read data from the range block.
func (s *SheetRange) Read() {}

// Write data to the range block.
func (s *SheetRange) Write() {}

// SetBold
func (s *SheetRange) SetBold(bold bool) {}

// SetItalic
func (s *SheetRange) SetItalic(set bool) {}

// SetFontSize
func (s *SheetRange) SetFontSize(set bool) {}

// SetTextDecoration
func (s *SheetRange) SetTextDecoration(set bool) {}

// SetHorizontalAlign
func (s *SheetRange) SetHorizontalAlign(set bool) {}

// SetVerticalAlign
func (s *SheetRange) SetVerticalAlign(set bool) {}

// SetFontColor
func (s *SheetRange) SetFontColor(set bool) {}

// SetBackgroudColor
func (s *SheetRange) SetBackgroudColor(set bool) {}

// SetBorder
func (s *SheetRange) SetBorder(borderType string, borderColor string) {}

// Clean all the style
func (s *SheetRange) Clean() {}

// Merge
func (s *SheetRange) Merge(mergeType string) {}

// Unmerge
func (s *SheetRange) Unmerge() {}

func (s *SheetRange) Find(keyword string, matchCase, matchEntireCell, regex, includeFormulas bool) {

}

func (s *SheetRange) Replace() {}
*/

// SetDropdown
// Parameter:
//  colors: set highlight color of values, could be nil. values like #1FB6C1
// Reference:
//  https://open.feishu.cn/document/ukTMukTMukTM/uATMzUjLwEzM14CMxMTN/datavalidation/set-dropdown
func (s *SheetRange) SetDropdown(values []string, multiple bool, colors []string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"range":              s.genRangeStr(),
		"dataValidationType": "list",
		"dataValidation": map[string]interface{}{
			"conditionValues": values,
			"options": map[string]interface{}{
				"multipleValues":     multiple,
				"highlightValidData": len(colors) > 0,
				"colors":             colors,
			},
		},
	})
	_url := s.sheet.ssClient.baseClient.urlJoin("open-apis/sheets/v2/spreadsheets/" + s.sheet.ssClient.GetToken() + "/dataValidation")
	req, _ := http.NewRequest(http.MethodPost, _url, bytes.NewReader(body))
	_, err := s.sheet.ssClient.baseClient.CommonReq(req, nil)
	return err
}

// Rows ...
func (s *SheetRange) Content() *SheetContent {
	sc := &SheetContent{}
	if s.Err != nil {
		sc.Err = s.Err
		return sc
	}
	content, err := s.sheet.getContentByRange(s.genRangeStr())
	if err != nil {
		sc.Err = err
		return sc
	}

	return content
}

func (s *SheetRange) genRangeStr() string {
	return s.sheet.id + "!" + s.leftTop.String() + ":" + s.rightBottom.String()
}
