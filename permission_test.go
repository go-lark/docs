package docs

import (
	"testing"

	"github.com/hilaily/kit/stringx"

	"github.com/stretchr/testify/assert"
)

func TestPermission_PublicSet(t *testing.T) {
	folderToken, err := getRootFolder()
	assert.NoError(t, err)
	r, err := getClientNew().file().Create(folderToken, "test create file", FileTypeDoc)
	t.Log(r)
	assert.NoError(t, err)
	err = getPermission().PublicSet(&PublicSet{
		Token:           r.Token,
		Type:            FileTypeDoc,
		LinkShareEntity: stringx.Pointer("tenant_readable"),
	})
	assert.NoError(t, err)
}

func TestPermission_Add(t *testing.T) {
	folderToken, err := getRootFolder()
	assert.NoError(t, err)
	r, err := getClientNew().file().Create(folderToken, "test create file", FileTypeDoc)
	assert.NoError(t, err)
	err = getPermission().Add(r.Token, FileTypeDoc, PermEdit, false, NewMemberWithEmail(testUserEmail))
	assert.NoError(t, err)
}

func TestPermission_TransferOwner(t *testing.T) {
	folderToken, err := getRootFolder()
	assert.NoError(t, err)
	r, err := getClientNew().file().Create(folderToken, "test create file", FileTypeDoc)
	assert.NoError(t, err)
	_, res, err := getPermission().TransferOwner(r.Token, FileTypeDoc, NewMemberWithEmail(testUserEmail), false, true)
	assert.NoError(t, err)
	assert.True(t, res.IsSuccess)
	t.Log(res)
}

func getPermission() *permission {
	return &permission{
		baseClient: getClientNew(),
	}
}

func getRootFolder() (string, error) {
	folder := getClientNew().RootFolder()
	if folder.Err != nil {
		return "", folder.Err
	}
	return folder.GetToken(), nil
}
