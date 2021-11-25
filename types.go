package docs

// FileType represent a document type of docs.
type FileType = string

var (
	FileTypeDoc      FileType = "doc"
	FileTypeSheet    FileType = "sheet"
	FileTypeSlide    FileType = "slide"
	FileTypeBitable  FileType = "bitable"
	FileTypeMindNote FileType = "mindnote"
	FileTypeFolder   FileType = "folder"
)

type ParentType = string

var (
	ParentTypeExplorer ParentType = "explorer"
)

type ModifySheetType = string

var (
	ModifySheetAdd    ModifySheetType = "addSheet"
	ModifySheetCopy   ModifySheetType = "copySheet"
	ModifySheetDelete ModifySheetType = "deleteSheet"
	ModifySheetUpdate ModifySheetType = "updateSheet"
)

type InseartDataOptionType = string

var (
	InseartDataOptionOverwrite  InseartDataOptionType = "OVERWRITE"
	InseartDataOptionInsertRows InseartDataOptionType = "INSERT_ROWS"
)

type MajorDimensionType = string

var (
	MajorDimensionRows    MajorDimensionType = "ROWS"
	MajorDimensionColumns MajorDimensionType = "COLUMNS"
)

type InheritStyleType = string

var (
	InheritStyleBefore InheritStyleType = "BEFORE"
	InheritStyleAfter  InheritStyleType = "AFTER"
)

type SheetRenderOption string

var (
	SheetRenderToString         SheetRenderOption = "ToString"
	SheetRenderFormattedValue   SheetRenderOption = "FormattedValue"
	SheetRenderFormula          SheetRenderOption = "Formula"
	SheetRenderUnformattedValue SheetRenderOption = "UnformattedValue"
)

type SheetDateTimeRenderOption string

var (
	SheetDateTimeRenderFormattedString SheetDateTimeRenderOption = "FormattedString"
)

type tokenIns struct {
	token string
}

func (t *tokenIns) GetToken() string {
	return t.token
}
