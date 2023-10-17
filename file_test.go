package docs

import (
	"bytes"
	"os"
	"testing"

	"github.com/go-lark/docs/doctypes"
	"github.com/stretchr/testify/assert"
)

func TestFile_Copy(t *testing.T) {
	folder := getClient().RootFolder()
	assert.NoError(t, folder.Err)
	_, res, err := getFile().Copy(testDocToken, FileTypeDoc, folder.GetToken(), "a copy file", true)
	assert.NoError(t, err)
	t.Logf("%#+v", res)
	err = getClient().permission().Add(res.Token, res.Type, PermEdit, false, NewMemberWithEmail(testUserEmail))
	assert.NoError(t, err)
}

func TestFile_Create(t *testing.T) {
	res, err := getFile().Create(testFolderToken, "test create file", FileTypeSheet)
	//res, err := getClientNew().File().Create(testFolderToken, "test create file1", FileTypeSheet)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.URL)
	t.Log(res.URL)
}

func TestFile_UpdateAll(t *testing.T) {
	f := &bytes.Buffer{}
	f.WriteString("test update file")
	res, err := getFile().UpdateAll(
		ParentTypeExplorer,
		testFolderToken,
		"testfile.txt",
		f.Bytes(),
	)
	assert.NoError(t, err)
	t.Log(res)
}

func TestFilePrepares(t *testing.T) {
	fi, err := os.Stat(testBigFile)
	assert.NoError(t, err)
	resp, err := getFile().resumePrepare(
		fileUpdateResumePrepare,
		fi.Name(),
		fi.Size(),
		doctypes.AttachmentFile,
		testDocToken,
	)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestFileUpdateResume(t *testing.T) {
	fi, err := os.Stat(testBigFile)
	assert.NoError(t, err)
	f, err := os.Open(testBigFile)
	assert.NoError(t, err)
	ch := make(chan int64, 20)
	go func() {
		for v := range ch {
			t.Log("process: ", v)
		}
	}()

	resp, err := getFile().UpdateResumed(
		ParentTypeExplorer,
		testFolderToken,
		fi.Name(),
		fi.Size(),
		f,
		ch,
	)
	assert.NoError(t, err)
	t.Log(resp)
}

func getFile() *File {
	return getClient().file()
}
