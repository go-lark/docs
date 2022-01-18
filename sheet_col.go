package docs

import "strings"

/*
	Method of sheet column.
*/

// cellname AA11, colname AA, rowname 11
func cellnameSplit(cellname string) (string, int) {
	b := []byte(cellname)
	i := 0
	for i = 0; i <= len(b)-1; i++ {
		if b[i] < 'A' || b[i] > 'Z' {
			break
		}
	}
	colnameByte := b[:i]
	rownameByte := b[i:]
	rowCount := 0
	base := 1
	for i := len(rownameByte) - 1; i >= 0; i-- {
		rowCount = (int(rownameByte[i]-'0'))*base + rowCount
		base *= 10
	}
	return string(colnameByte), rowCount
}

// 本质上是个 26 进制的计算
func colnameAdd(colname string, num int) string {
	i := colName2Num(colname)
	return num2ColName(i + num)
}

func num2ColName(num int) string {
	var col string
	for num > 0 {
		col = string(rune((num-1)%26+65)) + col
		num = (num - 1) / 26
	}
	return col
}

// A represent 1
func colName2Num(name string) int {
	name = strings.ToUpper(name)
	col := 0
	base := 1
	for i := len(name) - 1; i >= 0; i-- {
		r := name[i]
		col += int(r-'A'+1) * base
		base *= 26
	}
	return col
}
