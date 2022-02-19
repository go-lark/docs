package docs

import (
	"fmt"
	"reflect"
	"strconv"
)

// Range reference https://open.feishu.cn/document/ukTMukTMukTM/uczNzUjL3czM14yN3MTN#bae19f77
type SheetRange struct {
	Err         error
	leftTop     *cellName
	rightBottom *cellName
	sheet       *Sheet
}

/*
// Prepend insert data before range block.
func (s *SheetRange) Prepend() {}

// Append insert data after range block.
func (s *SheetRange) Append() {}

// Read data from the range block.
func (s *SheetRange) Read() {}

// Write data to the range block.
func (s *SheetRange) Write() {}

// SetBold
func (s *SheetRange) SetBold(bold bool) {}

// SetItalic
func (s *SheetRange) SetItalic(set bool) {}

// SetFontSize
func (s *SheetRange) SetFontSize(set bool) {}

// SetTextDecoration
func (s *SheetRange) SetTextDecoration(set bool) {}

// SetHorizontalAlign
func (s *SheetRange) SetHorizontalAlign(set bool) {}

// SetVerticalAlign
func (s *SheetRange) SetVerticalAlign(set bool) {}

// SetFontColor
func (s *SheetRange) SetFontColor(set bool) {}

// SetBackgroudColor
func (s *SheetRange) SetBackgroudColor(set bool) {}

// SetBorder
func (s *SheetRange) SetBorder(borderType string, borderColor string) {}

// Clean all the style
func (s *SheetRange) Clean() {}

// Merge
func (s *SheetRange) Merge(mergeType string) {}

// Unmerge
func (s *SheetRange) Unmerge() {}

func (s *SheetRange) Find(keyword string, matchCase, matchEntireCell, regex, includeFormulas bool) {

}

func (s *SheetRange) Replace() {}
*/

// Rows ...
func (s *SheetRange) Rows() ([]SheetRow, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	content, err := s.sheet.getContentByRange(s.genRangeStr())
	if err != nil {
		return nil, err
	}

	return content.ToRows(), nil
}

// RowsParseMerge
// parse cell from merged cell.
// if a cell is merged, like below, we only get value at the first cell, others are nil.
// so we should fill a to every cell of the merged cell.
// | a |   |   |
// |   |   |   |
// |   |   |   |
func (s *SheetRange) RowsParseMerge() ([]SheetRow, error) {
	meta, err := s.sheet.getMeta()
	if err != nil {
		return nil, err
	}
	cells, err := s.Rows()
	if err != nil {
		return nil, err
	}
	for _, v := range meta.Merges {
		startCol := v.StartColumnIndex
		endCol := v.StartColumnIndex + v.ColumnCount - 1
		startRow := v.StartRowIndex
		endRow := v.StartRowIndex + v.RowCount - 1
		for i := startRow; i <= endRow; i++ {
			for j := startCol; j <= endCol; j++ {
				cells[i][j] = cells[startRow][startCol]
			}
		}
	}
	return cells, nil
}

func (s *SheetRange) Scan(ptr interface{}) error {
	rows, err := s.RowsParseMerge()
	if err != nil {
		return err
	}
	return s.scan(rows, ptr)
}

func (s *SheetRange) scan(cells []SheetRow, ptr interface{}) error {
	// check it args is a pointer
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("ptr is must a slice pointer")
	}
	if rv.CanSet() {
		return fmt.Errorf("can not get slice address")
	}
	// dereference the pointer to get the slice
	rv = rv.Elem()
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("ptr element is must a slice")
	}
	// slice element
	elemt := rv.Type().Elem()
	if elemt.Kind() != reflect.Ptr {
		return fmt.Errorf("slice element must be a struct pointer, it is not pointer")
	}
	elemt = elemt.Elem()
	if elemt.Kind() != reflect.Struct {
		return fmt.Errorf("slice element must be a struct pointer, it it not struct")
	}

	fieldPosition := map[int]int{}
	for i := 0; i < elemt.NumField(); i++ {
		str := elemt.Field(i).Tag.Get("docs")
		if str != "" {
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return fmt.Errorf("struct tag wrong, it is not a number,  field index, %d, tag: %s", i, str)
			}
			fieldPosition[i] = int(v)
		} else {
			fieldPosition[i] = i
		}
	}
	for _, row := range cells {
		// new a slice element
		valP := reflect.New(elemt)
		// dereference pointer
		val := valP.Elem()
		for j := 0; j < val.NumField(); j++ {
			name := elemt.Field(j).Name
			position := fieldPosition[j]
			f := val.Field(j)
			if !f.CanAddr() {
				fmt.Printf(": can not set, %s\n", name)
			}
			switch f.Kind() {
			case reflect.Int64, reflect.Int:
				to, err := row[position].ToInt64()
				if err != nil {
					return fmt.Errorf("scan faild, index: %d, %w", j, err)
				}
				f.SetInt(to)
			case reflect.Float64, reflect.Float32:
				to, err := row[position].ToFloat()
				if err != nil {
					return fmt.Errorf("scan faild, index: %d, %w", j, err)
				}
				f.SetFloat(to)
			default:
				f.SetString(row[position].ToString())
			}

		}
		rv.Set(reflect.Append(rv, valP))
	}
	return nil
}

func (s *SheetRange) genRangeStr() string {
	return s.sheet.sheetMeta.SheetID + "!" + s.leftTop.String() + ":" + s.rightBottom.String()
}
