package docs

import (
	"bytes"
	"os"
	"testing"

	"github.com/go-lark/docs/doctypes"
	"github.com/stretchr/testify/assert"
)

func TestAttachment_UpdateAll(t *testing.T) {
	f := &bytes.Buffer{}
	f.WriteString("test update file")
	res, err := getAttachment().UpdateAll(
		doctypes.AttachmentFile,
		testDocToken,
		"testfile.txt",
		f.Bytes(),
	)
	assert.NoError(t, err)
	t.Log(res)
}

func TestAttachment_UpdateResume(t *testing.T) {
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

	resp, err := getAttachment().UpdateResuming(
		doctypes.AttachmentFile,
		testDocToken,
		fi.Name(),
		fi.Size(),
		f,
		ch,
	)
	assert.NoError(t, err)
	t.Log(resp)
}

func getAttachment() *Attachment {
	return getClient().attachment()
}
