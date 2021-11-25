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

type SheetCellStyle struct {
	Font struct {
		Bold     bool   `json:"bold,omitempty"`
		Italic   bool   `json:"italic,omitempty"`
		FontSize string `json:"font_size,omitempty"`
		Clean    bool   `json:"clean,omitempty"`
	} `json:"font,omitempty"`
	TextDecoration int    `json:"text_decoration,omitempty"`
	Formatter      string `json:"formatter,omitempty"`
	HAlign         int    `json:"h_align,omitempty"`
	VAlign         int    `json:"v_align,omitempty"`
	ForeColor      string `json:"fore_color,omitempty"`
	BackColor      string `json:"back_color,omitempty"`
	BorderType     string `json:"border_type,omitempty"`
	BorderColor    string `json:"border_color,omitempty"`
	Clean          bool   `json:"clean,omitempty"`
}
