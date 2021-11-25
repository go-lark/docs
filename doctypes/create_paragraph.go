package doctypes

import "github.com/go-lark/docs/log"

func NewBlockParagraph(elements ...IElement) *BlockParagraph {
	realElements := make([]*ParagraphElement, 0, len(elements))
	for _, v := range elements {
		realElements = append(realElements, v.ToElement())
	}
	p := &BlockParagraph{}
	p.Elements = realElements
	return p
}

func NewBlockParagraphWithTextRun(text string) *BlockParagraph {
	return NewBlockParagraph(NewElementTextRun(text))
}

func NewBlocksList(listType ListType, list []*IndentList) Blocks {
	blocks := make([]*Block, 0, len(list))
	for i, v := range list {
		if v.Ident < 1 || v.Ident > 16 {
			log.Errorln("[docs] list ident should in [1,16]")
		}
		b := NewBlockParagraph(NewElementTextRun(v.Text)).List(&List{
			Type:        listType,
			IndentLevel: v.Ident,
			Number:      i + 1,
		})
		blocks = append(blocks, b.ToBlocks()...)
	}
	return blocks
}

func NewBlockCode(codes []string) Blocks {
	list := make([]*IndentList, 0, len(codes))
	for _, v := range codes {
		list = append(list, &IndentList{Text: v, Ident: 0})
	}
	return NewBlocksList(ListCode, list)
}

func NewBlockCallout(callout *BlockCallout) *BlockCallout {
	if callout.CalloutTextColor != nil {
		callout.CalloutTextColor.Alpha = 1.0
	}
	if callout.CalloutBackgroundColor != nil {
		callout.CalloutBackgroundColor.Alpha = 1.0
	}
	if callout.CalloutBorderColor != nil {
		callout.CalloutBorderColor.Alpha = 1.0
	}
	return callout
}

func NewElementTextRun(text string) *TextRun {
	t := &TextRun{
		Text: text,
	}
	return t
}

func NewElementLink(url string) *Link {
	return &Link{
		URL: url,
	}
}

func NewElementDocsLink(url string) *DocsLink {
	return &DocsLink{
		URL: url,
	}
}

func NewElementPerson(openID string) *Person {
	return &Person{
		OpenID: openID,
	}
}

func NewElementReminder(isWholeDay bool, timestamp int64, shouldNotify bool, notifyType ReminderNotifyType) *Reminder {
	return &Reminder{
		IsWholeDay:   isWholeDay,
		Timestamp:    timestamp,
		ShouldNotify: shouldNotify,
		NotifyType:   notifyType,
	}
}

func NewElementEquation(equation string) *Equation {
	return &Equation{
		Equation: equation,
	}
}

func NewElementFileByToken(token string, viewType FileViewType) *ElementFile {
	return &ElementFile{
		FileToken: token,
		ViewType:  viewType,
	}
}
