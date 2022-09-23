package docs

// NOT SUPPORT NOW

/*
import (
	"io"

	"github.com/go-lark/docs/doctypes"
)
*/

/*
	for update a file or image in a doc
*/

/*
const (
	attachmentUpdateAllURL        = "/open-apis/drive/v1/medias/upload_all"
	attachmentUpdateResumePrepare = "/open-apis/drive/v1/medias/upload_prepare"
	attachmentUpdateResumePart    = "/open-apis/drive/v1/medias/upload_part"
	attachmentUpdateResumeFinish  = "/open-apis/drive/v1/medias/upload_finish"
)

func newAttachment(client *Client) *Attachment {
	return &Attachment{
		f: newFile(client),
	}
}

type Attachment struct {
	f *File
}

func (a *Attachment) UpdateAll(attachmentType doctypes.AttachmentType, token, filename string, fileData []byte) (string, error) {
	return a.f.updateAllBase(attachmentUpdateAllURL, attachmentType, token, filename, fileData)
}

func (a *Attachment) UpdateResuming(attachmentType doctypes.AttachmentType, token, filename string, fileSize int64, fileData io.Reader, processChan chan int64) (string, error) {
	return a.f.updateResumeBase(
		attachmentUpdateResumePrepare, attachmentUpdateResumePart, attachmentUpdateResumeFinish,
		attachmentType, token, filename,
		fileSize, fileData, processChan,
	)
}

func (a *Attachment) Delete(){

}
*/
