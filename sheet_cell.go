package docs

import "encoding/json"

/*
	Method of sheet cell.
*/

func NewSheetCell(i interface{}) *SheetCell {
	return &SheetCell{
		interData{val: i},
	}
}

// SheetCell represent for a cell of sheet
type SheetCell struct {
	interData
}

func (s SheetCell) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.interData.val)
}

/*
	data structure sheet support
	https://open.feishu.cn/document/ukTMukTMukTM/ugjN1UjL4YTN14CO2UTN
*/

// SheetCellTypeLink ...
func SheetCellTypeLink(title, link string) interface{} {
	return map[string]string{
		"text": title,
		"link": link,
		"type": "url",
	}
}

func SheetCellTypeMentionEmail(email string, notify, grantReadPermission bool) interface{} {
	return map[string]interface{}{
		"type":                "mention",
		"text":                email,
		"textType":            "email",
		"notify":              notify,
		"grantReadPermission": grantReadPermission,
	}
}

func SheetCellTypeMentionUnionID(unionID string, notify, grantReadPermission bool) interface{} {
	return map[string]interface{}{
		"type":                "mention",
		"text":                unionID,
		"textType":            "unionId",
		"notify":              notify,
		"grantReadPermission": grantReadPermission,
	}
}

func SheetCellTypeMentionOpenID(openID string, notify, grantReadPermission bool) interface{} {
	return map[string]interface{}{
		"type":                "mention",
		"text":                openID,
		"textType":            "openId",
		"notify":              notify,
		"grantReadPermission": grantReadPermission,
	}
}

func SheetCellTypeDocument(docType FileType, docToken string) interface{} {
	return map[string]string{
		"type":     "mention",
		"textType": "fileToken",
		"text":     docToken,
		"objType":  docType,
	}
}

func SheetCellTypeFormula(text string) interface{} {
	return map[string]string{
		"type": "formula",
		"text": text,
	}
}

func SheetCellTypeDropdown(text []interface{}) interface{} {
	return map[string]interface{}{
		"type":   "multipleValue",
		"values": text,
	}
}

type SheetCellStyle struct {
	Font struct {
		Bold     bool   `json:"bold,omitempty"`
		Italic   bool   `json:"italic,omitempty"`
		FontSize string `json:"fontSize,omitempty"`
		Clean    bool   `json:"clean,omitempty"`
	} `json:"font,omitempty"`
	TextDecoration int    `json:"textDecoration,omitempty"`
	Formatter      string `json:"formatter,omitempty"`
	HAlign         int    `json:"hAlign,omitempty"`
	VAlign         int    `json:"vAlign,omitempty"`
	ForeColor      string `json:"foreColor,omitempty"`
	BackColor      string `json:"backColor,omitempty"`
	BorderType     string `json:"borderType,omitempty"`
	BorderColor    string `json:"borderColor,omitempty"`
	Clean          bool   `json:"clean,omitempty"`
}
