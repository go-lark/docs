package docs

/*
import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-lark/docs/doctypes"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestAttachment_UpdateAll(t *testing.T) {
	Convey("TestAttachment_UpdateAll", t, func() {
		f := &bytes.Buffer{}
		f.WriteString("test update file")
		res, err := getAttachment().UpdateAll(
			doctypes.AttachmentFile,
			testDocToken,
			"testfile.txt",
			f.Bytes(),
		)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		assert.NoError(t, err)
		t.Log(res)
	})
}

func TestAttachment_UpdateResume(t *testing.T) {
	Convey("TestAttachment_UpdateResume", t, func() {
		path := os.TempDir()
		name := filepath.Join(path, "test.txt")
		t.Log(name)
		testBigFile, err := os.Create(name)
		So(err, ShouldBeNil)
		testBigFile.Truncate(2 * 1e7)
		fi, err := testBigFile.Stat()
		assert.NoError(t, err)
		testBigFile.Close()

		f, err := os.Open(name)
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
	})
}

func TestDeleteFile(t *testing.T) {
	res, err := getClientNew().file().delete("boxcnkaAlfia0qxcumSToGxLahe", FileTypeFile)
	assert.NoError(t, err)
	t.Log(res)
}

func getAttachment() *Attachment {
	return getClientNew().attachment()
}
*/
