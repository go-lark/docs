package doctypes

type Title struct {
	TitleStyle *ParagraphStyle     `json:"style,omitempty"`
	Elements   []*ParagraphElement `json:"elements"`
	LocationEmbed
}

func (t *Title) Align(align ParagraphAlignType) *Title {
	if t.TitleStyle == nil {
		t.TitleStyle = &ParagraphStyle{}
	}
	t.TitleStyle.Align = align
	return t
}

type Body struct {
	Blocks     []*Block    `json:"blocks,omitempty"`
	Attachment interface{} `json:"attachment,omitempty"`
}

type Location struct {
	ZoneID     string `json:"zoneID"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
}
