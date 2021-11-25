package doctypes

type IBlocks interface {
	//ToBlock() *Block
	ToBlocks() []*Block
}

// Block ...
type Block struct {
	Type           string               `json:"type"`
	Paragraph      *BlockParagraph      `json:"paragraph,omitempty"`
	HorizontalLine *BlockHorizontalLine `json:"horizontalLine,omitempty"`
	EmbeddedPage   *BlockEmbeddedPage   `json:"embeddedPage,omitempty"`
	ChatGroup      *BlockChatGroup      `json:"chatGroup,omitempty"`
	Table          *BlockTable          `json:"table,omitempty"`
	Sheet          *BlockSheet          `json:"sheet,omitempty"`
	Diagram        *BlockDiagram        `json:"diagram,omitempty"`
	Jira           *BlockJira           `json:"jira,omitempty"`
	Poll           *BlockPoll           `json:"poll,omitempty"`
	Bitable        *BlockBitable        `json:"bitable,omitempty"`
	UndefinedBlock *BlockUndefined      `json:"undefined_block,omitempty"`
	Gallery        *BlockGallery        `json:"gallery,omitempty"`
	Callout        *BlockCallout        `json:"callout,omitempty"`
	DocsApp        *BlockDocsApp        `json:"docsApp,omitempty"`
}

// BlockHorizontalLine  ...
type BlockHorizontalLine struct {
	LocationEmbed
}

func (h *BlockHorizontalLine) ToBlocks() []*Block {
	return []*Block{{
		Type:           blockHorizontalLine,
		HorizontalLine: h,
	}}
}

// BlockChatGroup ...
type BlockChatGroup struct {
	OpenChatID string `json:"openChatId"`
	LocationEmbed
}

func (c *BlockChatGroup) ToBlocks() []*Block {
	return []*Block{{
		Type:      blockChatGroup,
		ChatGroup: c,
	}}
}

// BlockSheet represent a sheet block
type BlockSheet struct {
	Token      string `json:"token"`
	RowSize    int    `json:"rowSize"`
	ColumnSize int    `json:"columnSize"`
	LocationEmbed
}

func (s *BlockSheet) ToBlocks() []*Block {
	return []*Block{{
		Type:  blockSheet,
		Sheet: s,
	}}
}

type BlockBitable struct {
	Token    string          `json:"token"`
	ViewType BitableViewType `json:"viewType"`
	LocationEmbed
}

func (b *BlockBitable) ToBlocks() []*Block {
	return []*Block{{
		Type:    blockBitable,
		Bitable: b,
	}}
}

type BlockDiagram struct {
	Token       string `json:"token"`
	DiagramType string `json:"diagramType"`
	LocationEmbed
}

type BlockJira struct {
	Token     string `json:"token"`
	BlockType string `json:"block_type"`
	LocationEmbed
}

type BlockPoll struct {
	Token string `json:"token"`
	LocationEmbed
}

type BlockUndefined struct {
	LocationEmbed
}

type LocationEmbed struct {
	Location *Location `json:"location,omitempty"`
}

type BlockGallery struct {
	GalleryStyle GalleryStyle `json:"galleryStyle"`
	ImageList    []ImageItem  `json:"imageList"`
	LocationEmbed
}

type GalleryStyle struct {
	Align string `json:"align"`
}

type ImageItem struct {
	FileToken string `json:"fileToken"`
}

func (g *BlockGallery) ToBlocks() []*Block {
	return []*Block{{
		Type:    blockGallery,
		Gallery: g,
	}}
}

// BlockCallout ...
type (
	BlockCallout struct {
		CalloutEmojiID         string           `json:"calloutEmojiId"`
		CalloutBackgroundColor *ColorRGBA       `json:"calloutBackgroundColor"`
		CalloutBorderColor     *ColorRGBA       `json:"calloutBorderColor"`
		CalloutTextColor       *ColorRGBA       `json:"calloutTextColor"`
		Body                   BlockCalloutBody `json:"body"`
		ZoneId                 string           `json:"zoneId"`
		LocationEmbed
	}
	BlockCalloutBody struct {
		Blocks []*Block `json:"blocks"`
	}
)

func (b *BlockCallout) ToBlocks() []*Block {
	return []*Block{{
		Type:    "callout",
		Callout: b,
	}}
}

// BlockDocsApp ...
type BlockDocsApp struct {
	TypeID     string `json:"typeId"`
	InstanceID string `json:"instanceID"`
}

func (b *BlockDocsApp) ToBlocks() []*Block {
	return []*Block{{
		Type:    "docsApp",
		DocsApp: b,
	}}
}
