package doctypes

import (
	"net/url"
	"strings"

	"github.com/go-lark/docs/log"
)

// NewTitle create a doc title
func NewTitle(content string) *Title {
	t := &Title{}
	t.Elements = []*ParagraphElement{
		{Type: elementTextRun, TextRun: NewElementTextRun(content)},
	}
	return t
}

// NewBody create a doc body
func NewBody(blocks ...IBlocks) *Body {
	realBlocks := make([]*Block, 0, len(blocks))
	for _, v := range blocks {
		realBlocks = append(realBlocks, v.ToBlocks()...)
	}
	b := &Body{
		Blocks: realBlocks,
	}
	return b
}

// NewBlockHorizontalLine ...
func NewBlockHorizontalLine() *BlockHorizontalLine {
	return &BlockHorizontalLine{}
}

// NewBlockChatGroup
// @param: openChatID represent a open id for a chat, normally start with 'oc'
func NewBlockChatGroup(openChatID string) *BlockChatGroup {
	return &BlockChatGroup{
		OpenChatID: openChatID,
	}
}

// NewBlockSheet represent create a sheet block
// @params: token represent a sheet docs token.
// NOTE: if token is empty, create a new sheet block, if token is not empty create a sheet block copy from the existed sheet.
func NewBlockSheet(token string, rowSize, colSize int) *BlockSheet {
	if rowSize > 9 || colSize > 9 {
		log.Errorln("new block sheet, row size or col size is greater than 9, see: https://open.feishu.cn/document/ukTMukTMukTM/ukDM2YjL5AjN24SOwYjN#53fe05b8")
	}
	return &BlockSheet{
		Token:      token,
		RowSize:    rowSize,
		ColumnSize: colSize,
	}
}

// NewBlockBitable for new bitable block
// @param: token represent a bitable docs token
// NOTE:
// if token is empty create a new bitable block, if token is not empty create a bitable copy from the existed bitable.
func NewBlockBitable(token string, viewType BitableViewType) *BlockBitable {
	return &BlockBitable{
		Token:    token,
		ViewType: viewType,
	}
}

// NewBlockTable represent create a table block
func NewBlockTable(rowSize, colSize int) *BlockTable {
	return &BlockTable{
		RowSize:    rowSize,
		ColumnSize: colSize,
	}
}

// NewBlockEmbeddedPage ...
func NewBlockEmbeddedPage(embedType EmbedType, pageurl string) *BlockEmbeddedPage {
	if embedType == EmbedBilibili {
		if !strings.Contains(pageurl, "player.bilibili.com") {
			log.Errorln("embed bilibili page should use url in iframe, common start with palyer.bilibili.com")
		}
	}
	return &BlockEmbeddedPage{
		Type: embedType,
		URL:  escapeURL(pageurl),
	}
}

func escapeURL(u string) string {
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return url.QueryEscape(u)
	}
	return u
}
