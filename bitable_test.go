package docs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	c := bitableClient()
	err := c.Table("tblzYfTu6y2UVUIR").AddRecord([]Field{
		{
			"多行文本": "多行文本",
			"单选":   "单选",
			"日期":   time.Now().UnixMilli(),
		},
	}).Err
	assert.NoError(t, err)
}

func bitableClient() *Bitable {
	return getClientNew().OpenBitable("FYsxbceSuauxgzsmFVpc1fFEnh1")
}
