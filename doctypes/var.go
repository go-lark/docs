package doctypes

const (
	blockParagraph      = "paragraph"
	blockHorizontalLine = "horizontalLine"
	blockEmbeddedPage   = "embeddedPage"
	blockChatGroup      = "chatGroup"
	blockTable          = "table"
	blockSheet          = "sheet"
	blockBitable        = "bitable"
	blockGallery        = "gallery"
	//blockFile           = "file"
	//blockDiagram        = "diagram"
	//blockJira           = "jira"
	//blockPoll           = "poll"
	//blockUndefined      = "undefinedBlock"

	elementTextRun  = "textRun"
	elementDocsLink = "docsLink"
	elementPerson   = "person"
	elementEquation = "equation"
	elementFile     = "file"
	elementReminder = "reminder"
	//elementUndefined = "undefinedElement"
)

type AttachmentType = string

var (
	AttachmentImage AttachmentType = "doc_image"
	AttachmentFile  AttachmentType = "doc_file"
)

// ParagraphHeadineLevel represent paragraph style heading level
// Doc: https://open.feishu.cn/document/ukTMukTMukTM/ukDM2YjL5AjN24SOwYjN#4b468696
type ParagraphHeadineLevel = int

var (
	ParagraphHeadineLevel1 ParagraphHeadineLevel = 1
	ParagraphHeadineLevel2 ParagraphHeadineLevel = 2
	ParagraphHeadineLevel3 ParagraphHeadineLevel = 3
	ParagraphHeadineLevel4 ParagraphHeadineLevel = 4
	ParagraphHeadineLevel5 ParagraphHeadineLevel = 5
	ParagraphHeadineLevel6 ParagraphHeadineLevel = 6
	ParagraphHeadineLevel7 ParagraphHeadineLevel = 7
	ParagraphHeadineLevel8 ParagraphHeadineLevel = 8
	ParagraphHeadineLevel9 ParagraphHeadineLevel = 9
)

type ListType = string

var (
	ListNumber     = "number"
	ListBullet     = "bullet"
	ListCheckBox   = "checkBox"
	ListCheckedBox = "checkedBox"
	ListCode       = "code"
)

type FileViewType = string

var (
	FileViewPreview FileViewType = "preview"
	FileViewCard    FileViewType = "card"
	FileViewInline  FileViewType = "inline"
)

type ReminderNotifyType = int

var (
	ReminderNotifyNow       ReminderNotifyType = 0
	ReminderNotify5mBefore  ReminderNotifyType = 1
	ReminderNotify15mBefore ReminderNotifyType = 2
	ReminderNotify30mBefore ReminderNotifyType = 3
	ReminderNotify1hBefore  ReminderNotifyType = 4
	ReminderNotify2hBefore  ReminderNotifyType = 5
	ReminderNotify1dBefore  ReminderNotifyType = 6
	ReminderNotify2dBefore  ReminderNotifyType = 7
	// TODO Laily more types
)

type ParagraphAlignType = string

var (
	ParagraphAlignLeft   ParagraphAlignType = "left"
	ParagraphAlignRight  ParagraphAlignType = "right"
	ParagraphAlignCenter ParagraphAlignType = "center"
)

type EmbedType = string

var (
	EmbedBilibili EmbedType = "bilibili"
	EmbedXigua    EmbedType = "xigua"
	EmbedYouku    EmbedType = "youku"
	EmbedAirtable EmbedType = "airtable"
	EmbedBaiduMap EmbedType = "baidumap"
)

type BitableViewType = string

var (
	BitableViewGrid   BitableViewType = "grid"
	BitableViewKanban BitableViewType = "kanban"
)
