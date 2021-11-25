package docs

type SheetRange struct {
}

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
