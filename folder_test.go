package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderRoot(t *testing.T) {
	res := getClientNew().RootFolder()
	t.Log(res.GetToken())
}

func TestFolder(t *testing.T) {
	title := "create_folder_test"
	res := getFolder().CreateSubFolder(title)
	assert.NoError(t, res.Err)

	subFolder := getClientNew().OpenFolder(res.GetToken())
	meta, err := subFolder.GetMeta()
	assert.NoError(t, err)
	pMeta, err := getFolder().GetMeta()
	assert.NoError(t, err)
	assert.Equal(t, title, meta.Name)
	assert.Equal(t, pMeta.ID, meta.ParentID)

	list, err := getFolder().Children(nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list.Children), 1)

	has := false
	for _, v := range list.Children {
		if v.Name == title {
			has = true
		}
	}
	assert.True(t, has)
}

func TestFolder_Children(t *testing.T) {
	res, err := getFolder().Children([]FileType{FileTypeDoc})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res.Children))

}

func getFolder() *Folder {
	return getClientNew().OpenFolder(testFolderToken)
}
