package docs

import (
	"testing"

	"github.com/hilaily/kit/helper"
	"github.com/stretchr/testify/assert"
)

func TestGetDocxMeta(t *testing.T) {
	res, err := docxClient().GetMeta()
	assert.NoError(t, err)
	assert.Equal(t, baseDomain+"/docx/"+testDocxToken, res.URL)
	assert.NotEmpty(t, res.LatestModifyUser)
	helper.PJSON(res)
}

func docxClient() *Docx {
	return getClientNew().OpenDocx(testDocxToken)
}
