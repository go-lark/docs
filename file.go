package docs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hilaily/kit/httpx"
)

/*
	for update a file
	reference：https://open.feishu.cn/document/ukTMukTMukTM/uUjM5YjL1ITO24SNykjN
*/

const (
	fileUpdateAllURL        = "/open-apis/drive/v1/files/upload_all"
	fileUpdateResumePrepare = "/open-apis/drive/v1/files/upload_prepare"
	fileUpdateResumePart    = "/open-apis/drive/v1/files/upload_part"
	fileUpdateResumeFinish  = "/open-apis/drive/v1/files/upload_finish"

	retryTimes = 3
)

func newFile(client *Client) *File {
	return &File{baseClient: client}
}

type File struct {
	baseClient *Client
}

// Create a doc or sheet
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uUTN5UjL1UTO14SN1kTN
// 创建一个文件
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uQTNzUjL0UzM14CN1MTN
func (f *File) Create(folderToken, title string, fileType FileType) (*RespCreateFile, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"title": title,
		"type":  fileType,
	})
	req, _ := http.NewRequest(http.MethodPost, f.baseClient.urlJoin("/open-apis/drive/explorer/v2/file/"+folderToken), bytes.NewReader(en))
	resp := &RespCreateFile{}
	_, err := f.baseClient.CommonReq(req, resp)
	return resp, err
}

// Copy a file
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uYTN5UjL2UTO14iN1kTN
// 复制一个文档
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uYTNzUjL2UzM14iN1MTN
func (f *File) Copy(srcFileToken string, srcFileType FileType, dstFolderToken, dstTitle string, copyComment bool) ([]byte, *RespCopyFile, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"type":           srcFileType,
		"dstFolderToken": dstFolderToken,
		"dstName":        dstTitle,
		"commentNeeded":  copyComment,
	})
	req, _ := http.NewRequest(http.MethodPost, f.baseClient.urlJoin("/open-apis/drive/explorer/v2/file/copy/files/"+srcFileToken), bytes.NewReader(en))
	resp := &RespCopyFile{}
	b, err := f.baseClient.CommonReq(req, resp)
	return b, resp, err
}

// UpdateAll
// Return
//
//	1: token of the file
func (f *File) UpdateAll(parentType ParentType, parentNode, filename string, fileData []byte) (string, error) {
	return f.updateAllBase(fileUpdateAllURL, parentType, parentNode, filename, fileData)
}

func (f *File) UpdateResumed(parentType ParentType, parentNode, filename string, fileSize int64, fileData io.Reader, processChan chan int64) (string, error) {
	return f.updateResumeBase(
		fileUpdateResumePrepare, fileUpdateResumePart, fileUpdateResumeFinish,
		parentType, parentNode, filename,
		fileSize, fileData, processChan,
	)
}

type UpdateFileResp struct {
	FileToken string `json:"file_token"`
}

func (f *File) updateAllBase(updateURL string, parentType, token, filename string, fileData []byte) (string, error) {
	checksum := adler32.Checksum(fileData)
	ct, body, err := httpx.NewFormBody(map[string]string{
		"parent_type": parentType,
		"parent_node": token,
		"file_name":   filename,
		"size":        strconv.Itoa(len(fileData)),
		"checksum":    strconv.Itoa(int(checksum)),
	},
		[]*httpx.FileInfo{
			{
				Fieldname: "file",
				Filename:  filename,
				Data:      bytes.NewBuffer(fileData),
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("gen http form body, %w", err)
	}
	_req, err := http.NewRequest(
		http.MethodPost,
		f.baseClient.urlJoin(updateURL),
		body,
	)
	if err != nil {
		return "", fmt.Errorf("gen http request, %w", err)
	}
	_req.Header.Set("Content-Type", ct)
	dst := &UpdateFileResp{}
	_, err = f.baseClient.DoRequest(_req, dst)
	if err != nil {
		return "", fmt.Errorf("http request, %w", err)
	}
	return dst.FileToken, nil
}

func (f *File) updateResumeBase(prepareURL, partURL, finishURL string, parentType ParentType, parentNode, filename string, fileSize int64, fileData io.Reader, processChan chan int64) (fileToken string, err error) {
	resumeInfo, err := f.resumePrepare(prepareURL, filename, fileSize, parentType, parentNode)
	if err != nil {
		return
	}
	if x, ok := fileData.(io.ReadCloser); ok {
		defer x.Close()
	}
	buf := make([]byte, resumeInfo.BlockSize)
	seq := 0
	process := int64(0)
	var n int
	ch := make(chan int64, int64(fileSize/resumeInfo.BlockSize)+1)
	if processChan != nil {
		go func() {
			for v := range ch {
				processChan <- v
			}
		}()
	}
	for {
		n, err = fileData.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			err = fmt.Errorf("update resume, read data, %w", err)
			return
		}
		if n == 0 {
			break
		}

		err = f.resumePart(partURL, resumeInfo.UploadID, seq, n, buf[:n])
		if err != nil {
			err = fmt.Errorf("update resume, part, %w", err)
			return
		}
		process += int64(n)
		ch <- process
		seq++
		//buf = buf[:0]
	}
	close(ch)
	fileToken, err = f.resumeFinish(finishURL, resumeInfo.UploadID, resumeInfo.BlockNum)
	return
}

func (f *File) resumePrepare(urlpath, filename string, fileSize int64, parentType string, parentNode string) (*resumePreparesResp, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"file_name":   filename,
		"size":        fileSize,
		"parent_type": parentType,
		"parent_node": parentNode,
	})
	req, err := http.NewRequest(http.MethodPost, f.baseClient.urlJoin(urlpath), bytes.NewReader(en))
	if err != nil {
		return nil, fmt.Errorf("gen resume prepare request, %w", err)
	}
	dst := &resumePreparesResp{}
	_, err = f.baseClient.CommonReq(req, dst)
	if err != nil {
		return nil, fmt.Errorf("resume prepares request, %w", err)
	}
	return dst, nil
}

func (f *File) resumePartReTry(urlpath, uploadID string, seq int, size int, partData []byte) (err error) {
	duration := time.Second
	for i := 1; i <= retryTimes; i++ {
		err = f.resumePart(urlpath, uploadID, seq, size, partData)
		if err != nil {
			return nil
		}
		time.Sleep(duration * time.Duration(i))
	}
	return err
}

func (f *File) resumePart(urlpath, uploadID string, seq int, size int, partData []byte) (err error) {
	checksum := adler32.Checksum(partData)
	ct, body, err := httpx.NewFormBody(
		map[string]string{
			"upload_id": uploadID,
			"seq":       strconv.Itoa(seq),
			"size":      strconv.Itoa(size),
			"checksum":  strconv.Itoa(int(checksum)),
		},
		[]*httpx.FileInfo{
			{
				Fieldname: "file",
				Data:      bytes.NewReader(partData),
			},
		},
	)
	if err != nil {
		err = fmt.Errorf("resume part, gen http request, upload_id: %s, seq: %d, %w", uploadID, seq, err)
		return
	}
	req, err := http.NewRequest(
		http.MethodPost,
		f.baseClient.urlJoin(urlpath),
		body,
	)
	if err != nil {
		err = fmt.Errorf("resume part, gen http request, %w", err)
		return
	}
	req.Header.Set("Content-Type", ct)
	_, err = f.baseClient.DoRequest(req, nil)
	if err != nil {
		err = fmt.Errorf("resume part, upload_id: %s, seq: %d, %w", uploadID, seq, err)
		return err
	}
	return
}

func (f *File) resumeFinish(urlpath string, uploadID string, blockSum int) (fileToken string, err error) {
	en, _ := json.Marshal(map[string]interface{}{
		"upload_id": uploadID,
		"block_num": blockSum,
	})
	req, err := http.NewRequest(http.MethodPost, f.baseClient.urlJoin(urlpath), bytes.NewReader(en))
	if err != nil {
		err = fmt.Errorf("resume finiesh, upload_id: %s, %w", uploadID, err)
		return
	}
	dst := &UpdateFileResp{}
	_, err = f.baseClient.CommonReq(req, dst)
	if err != nil {
		err = fmt.Errorf("resume finish, upload_id: %s, %w", uploadID, err)
		return
	}
	fileToken = dst.FileToken
	return
}

func (f *File) statistics(fileToken string, fileType FileType) (stat *FileStatistics, err error) {
	u := f.baseClient.urlJoin(fmt.Sprintf("open-apis/drive/v1/files/%s/statistics?file_type=%s", fileToken, fileType))
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		err = fmt.Errorf("gen request fail, token:%s, file type:%s, %w", fileToken, fileType, err)
		return
	}
	_, err = f.baseClient.CommonReq(req, &stat)
	if err != nil {
		err = fmt.Errorf("get statistics fail, token:%s, file type:%s, %w", fileToken, fileType, err)
		return
	}
	return
}

type resumePreparesResp struct {
	UploadID  string `json:"upload_id"`
	BlockSize int64  `json:"block_size"`
	BlockNum  int    `json:"block_num"`
}

type FileStatistics struct {
	FileToken  string `json:"file_token"`
	FileType   string `json:"file_type"`
	Statistics Stats  `json:"statistics"`
}

type Stats struct {
	UV        int64 `json:"uv"`
	PV        int64 `json:"pv"`
	LikeCount int64 `json:"like_count"`
	Timestamp int64 `json:"timestamp"`
}
