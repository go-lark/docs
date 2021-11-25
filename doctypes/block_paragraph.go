package doctypes

// Paragraph represent a doc paragraph
type BlockParagraph struct {
	ParagraphStyle ParagraphStyle      `json:"style"`    // 段落样式
	Elements       []*ParagraphElement `json:"elements"` // 段落元素
	LocationEmbed
}

func (p *BlockParagraph) ToBlocks() []*Block {
	return []*Block{{
		Type:      blockParagraph,
		Paragraph: p,
	}}
}

// HeadingLevel represent heading level style
func (p *BlockParagraph) HeadingLevel(level ParagraphHeadineLevel) *BlockParagraph {
	p.ParagraphStyle.HeadingLevel = level
	return p
}

// List ...
func (p *BlockParagraph) List(l *List) *BlockParagraph {
	p.ParagraphStyle.List = l
	return p
}

func (p *BlockParagraph) SetQuote() *BlockParagraph {
	p.ParagraphStyle.Quote = true
	return p
}

func (p *BlockParagraph) Align(align ParagraphAlignType) *BlockParagraph {
	p.ParagraphStyle.Align = align
	return p
}

type ParagraphStyle struct {
	HeadingLevel ParagraphHeadineLevel `json:"headingLevel,omitempty"`
	List         *List                 `json:"list"`
	Quote        bool                  `json:"quote"`
	Align        ParagraphAlignType    `json:"align"`
}

// TODO: Laily
type List struct {
	Type        string `json:"type"`
	IndentLevel int    `json:"indentLevel"`
	Number      int    `json:"number"`
}

type IElement interface {
	ToElement() *ParagraphElement
}

// ParagraphElement ...
type ParagraphElement struct {
	Type      string       `json:"type"`
	TextRun   *TextRun     `json:"textRun,omitempty"`
	DocsLink  *DocsLink    `json:"docsLink,omitempty"`
	Person    *Person      `json:"person,omitempty"`
	Equation  *Equation    `json:"equation,omitempty"`
	Reminider *Reminder    `json:"reminder,omitempty"`
	File      *ElementFile `json:"file,omitempty"`
}

type TextRun struct {
	Text      string     `json:"text"`
	TextStyle *TextStyle `json:"style,omitempty"`
	LineID    string     `json:"lineID"`
	LocationEmbed
}

func (t *TextRun) ToElement() *ParagraphElement {
	return &ParagraphElement{
		Type:    elementTextRun,
		TextRun: t,
	}
}

func (t *TextRun) SetBold() *TextRun {
	t.setStyle().Bold = true
	return t
}

func (t *TextRun) SetItalic() *TextRun {
	t.setStyle().Italic = true
	return t
}

func (t *TextRun) SetStrickThrouth() *TextRun {
	t.setStyle().StrikeThrough = true
	return t
}

func (t *TextRun) SetUnderline() *TextRun {
	t.setStyle().Underline = true
	return t
}

func (t *TextRun) SetCodeInline() *TextRun {
	t.setStyle().CodeInline = true
	return t
}

func (t *TextRun) SetColor(backgroudColor, textColor *ColorRGBA) *TextRun {
	t.setStyle().BackColor = backgroudColor
	t.setStyle().TextColor = textColor
	return t
}

func (t *TextRun) SetLink(link string) *TextRun {
	t.setStyle().Link = &Link{URL: link}
	return t
}

type TextStyle struct {
	Bold          bool       `json:"bold"`
	Italic        bool       `json:"italic"`
	StrikeThrough bool       `json:"strikeThrough"`
	Underline     bool       `json:"underline"`
	CodeInline    bool       `json:"codeInline"`
	BackColor     *ColorRGBA `json:"backColor,omitempty"`
	TextColor     *ColorRGBA `json:"textColor,omitempty"`
	Link          *Link      `json:"link"`
}

type Link struct {
	URL string `json:"url"`
}

type DocsLink struct {
	URL string `json:"url"`
	LocationEmbed
}

func (d *DocsLink) ToElement() *ParagraphElement {
	return &ParagraphElement{
		Type:     elementDocsLink,
		DocsLink: d,
	}
}

type Person struct {
	OpenID string `json:"openId"`
	LocationEmbed
}

func (p *Person) ToElement() *ParagraphElement {
	return &ParagraphElement{
		Type:   elementPerson,
		Person: p,
	}
}

type Reminder struct {
	IsWholeDay   bool  `json:"isWholeDay"`
	Timestamp    int64 `json:"timestamp"`
	ShouldNotify bool  `json:"shouldNotify"`
	NotifyType   int   `json:"notifyType"`
	LocationEmbed
}

func (r *Reminder) ToElement() *ParagraphElement {
	return &ParagraphElement{
		Type:      elementReminder,
		Reminider: r,
	}
}

type Equation struct {
	Equation string `json:"equation"`
	LocationEmbed
}

func (e *Equation) ToElement() *ParagraphElement {
	return &ParagraphElement{
		Type:     elementEquation,
		Equation: e,
	}
}

type ElementFile struct {
	FileToken string       `json:"fileToken"`
	ViewType  FileViewType `json:"viewType"`
	FileName  string       `json:"fileName"`
	Location  *Location    `json:"location,omitempty"`
}

func (f *ElementFile) ToElement() *ParagraphElement {
	return &ParagraphElement{
		Type: elementFile,
		File: f,
	}
}

// type custom

type IndentList struct {
	Text  string
	Ident int // [1,16]
}

type Blocks []*Block

func (b Blocks) ToBlocks() []*Block {
	return b
}

func (t *TextRun) setStyle() *TextStyle {
	if t.TextStyle == nil {
		t.TextStyle = &TextStyle{}
	}
	return t.TextStyle
}
