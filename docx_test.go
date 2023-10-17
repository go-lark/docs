package docs

import (
	"testing"

	"github.com/hilaily/kit/helper"
	"github.com/stretchr/testify/assert"
)

func TestGetDocxMeta(t *testing.T) {
	res, err := docxClient().GetMeta()
	assert.NoError(t, err)
	assert.Equal(t, tenantDomain+"/docx/"+testDocxToken, res.URL)
	assert.NotEmpty(t, res.LatestModifyUser)
	helper.PJSON(res)
}

func docxClient() *Docx {
	return getClient().OpenDocx(testDocxToken)
}
