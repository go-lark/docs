package doctypes

// BlockTable ...
type BlockTable struct {
	TableID    string       `json:"tableId"`
	RowSize    int          `json:"rowSize"`
	ColumnSize int          `json:"columnSize"`
	TableRows  []*TableRow  `json:"tableRows"`
	TableStyle *TableStyle  `json:"tableStyle"`
	MergeCells []*MergeCell `json:"mergeCells"`
	LocationEmbed
}

func (t *BlockTable) ToBlocks() []*Block {
	return []*Block{{
		Type:  blockTable,
		Table: t,
	}}
}

// Rows represent add rows for table
// []*TableRow create by NewTableRows
func (t *BlockTable) Rows(rows []*TableRow) *BlockTable {
	t.TableRows = rows
	return t
}

// StyleWith represent set width of columns
func (t *BlockTable) StyleWith(width int) *BlockTable {
	if t.TableStyle == nil {
		t.TableStyle = &TableStyle{}
	}
	p := t.TableStyle.TableColumnProperties
	p = append(p, TableColumnProperties{Width: width})
	t.TableStyle.TableColumnProperties = p
	return t
}

// MergedCell represent create a merged cell, can call the function many times to create more merged cell
func (t *BlockTable) MergedCell(rowStartIndex, rowEndIndex, colStartIndex, colEndIndex int) *BlockTable {
	cell := &MergeCell{
		RowStartIndex:    rowStartIndex,
		RowEndIndex:      rowEndIndex,
		ColumnStartIndex: colStartIndex,
		ColumnEndIndex:   colEndIndex,
	}
	t.MergeCells = append(t.MergeCells, cell)
	return t
}

// NewTableRows represent create some rows by cell content
func NewTableRows(cellContents [][]interface{}) []*TableRow {
	rows := make([]*TableRow, 0, len(cellContents))
	for ri, v := range cellContents {
		cells := make([]*TableCell, 0, len(v))
		for k, c := range v {
			cells = append(cells, &TableCell{
				ColumnIndex: k,
				Body:        c,
			})
		}
		rows = append(rows, &TableRow{
			RowIndex:   ri,
			TableCells: cells,
		})
	}
	return rows
}

type TableRow struct {
	RowIndex   int          `json:"rowIndex"`
	TableCells []*TableCell `json:"tableCells"`
}

// TableCell represent a normal table cell
type TableCell struct {
	ColumnIndex int         `json:"columnIndex"`
	ZoneID      string      `json:"zoneID"`
	Body        interface{} `json:"body"`
}

type TableStyle struct {
	TableColumnProperties []TableColumnProperties `json:"tableColumnProperties"`
}

type TableColumnProperties struct {
	Width int `json:"width"`
}

// MergeCell  represent a merged cell
type MergeCell struct {
	MergeCellID      string `json:"mergeCellId"`
	RowStartIndex    int    `json:"rowStartIndex"`
	RowEndIndex      int    `json:"rowEndIndex"`
	ColumnStartIndex int    `json:"columnStartIndex"`
	ColumnEndIndex   int    `json:"columnEndIndex"`
}
