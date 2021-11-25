package doctypes

// BlockEmbeddedPage represent a block which embedded a website page
type BlockEmbeddedPage struct {
	Type   EmbedType `json:"type"`
	URL    string    `json:"url"`
	Width  float64   `json:"width"`
	Height float64   `json:"height"`
	LocationEmbed
}

func (e *BlockEmbeddedPage) ToBlocks() []*Block {
	return []*Block{{
		Type:         blockEmbeddedPage,
		EmbeddedPage: e,
	}}
}

// WidthAndHeight for set width or height for the embedded website page
func (e *BlockEmbeddedPage) WidthAndHeight(width, height float64) *BlockEmbeddedPage {
	e.Width = width
	e.Height = height
	return e
}
