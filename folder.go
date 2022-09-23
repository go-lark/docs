package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-lark/docs/doctypes"
)

func newFolder(token string, client *Client) *Folder {
	f := &Folder{baseClient: client}
	f.token = token
	return f
}

// Folder represent a folder instance
type Folder struct {
	Err        error   // used for save error information
	tokenIns           // save token of the folder
	baseClient *Client // client instance of the app
}

// CreateSubFolder to create a child folder in the current folder
// Return
//
//	1: the child folder instance
func (f *Folder) CreateSubFolder(title string) *Folder {
	if f.Err != nil {
		return f
	}
	u := fmt.Sprintf("%s/open-apis/drive/explorer/v2/folder/%s", f.baseClient.domain, f.token)
	en, _ := json.Marshal(map[string]string{
		"title": title,
	})
	_req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(en))
	res := &RespCreateFoler{}
	_, err := f.baseClient.CommonReq(_req, res)
	if err != nil {
		f.Err = newErr(err.Error())
		return f
	}
	return newFolder(res.Token, f.baseClient)
}

// GetMeta to get the meta information of folder
func (f *Folder) GetMeta() (*FolderMeta, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	u := fmt.Sprintf("%s/open-apis/drive/explorer/v2/folder/%s/meta", f.baseClient.domain, f.token)
	_req, _ := http.NewRequest(http.MethodGet, u, nil)
	res := &FolderMeta{}
	_, err := f.baseClient.CommonReq(_req, res)
	return res, err
}

func (f *Folder) Children(fileTypes []FileType) (*ChildrenInfo, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	params := url.Values{}
	for _, v := range fileTypes {
		params.Add("types", v)
	}
	u, _ := url.Parse(fmt.Sprintf("%s/open-apis/drive/explorer/v2/folder/%s/children", f.baseClient.domain, f.token))
	u.RawQuery = params.Encode()
	_req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	res := &ChildrenInfo{}
	_, err := f.baseClient.CommonReq(_req, res)
	return res, err
}

// FIXME: Laily
// CreateDoc for create a document
func (f *Folder) CreateDoc(title *doctypes.Title, body *doctypes.Body) (doc *Doc) {
	doc = &Doc{}
	if f.Err != nil {
		doc.Err = f.Err
		return doc
	}
	content, _ := json.Marshal(&DocContent{
		Title: title,
		Body:  body,
	})
	postBody, _ := json.Marshal(map[string]interface{}{
		"FolderToken": f.token,
		"Content":     string(content),
	})
	u := fmt.Sprintf("%s/open-apis/doc/v2/create", f.baseClient.domain)
	_req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(postBody))
	res := &RespCreateDoc{}
	_, err := f.baseClient.CommonReq(_req, res)
	if err != nil {
		doc.Err = newErr(err.Error())
		return
	}
	return newDoc(res.ObjToken, f.baseClient)
}

// CreateSpreadSheet ...
func (f *Folder) CreateSpreadSheet(title string) (ss *SpreadSheets) {
	ss = &SpreadSheets{}
	if f.Err != nil {
		ss.Err = f.Err
		return
	}
	res, err := f.baseClient.file().Create(f.token, title, FileTypeSheet)
	if err != nil {
		ss.Err = err
		return
	}
	return newSpreadSheets(res.Token, f.baseClient)
}

// UploadFile is used for uploading small file(less than 20MB)
func (f *Folder) UploadFile(filename string, fileData []byte) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.baseClient.file().UpdateAll(ParentTypeExplorer, f.token, filename, fileData)
}

// UpdateFileResumed is used for uploading big file
func (f *Folder) UpdateFileResumed(filename string, fileSize int64, fileData io.Reader, processChan chan int64) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.baseClient.file().UpdateResumed(ParentTypeExplorer, f.token, filename, fileSize, fileData, processChan)
}

// RootFolder get root folder meta
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/ucDM3UjL3AzN14yNwcTN
// 获取我的空间 meta 信息
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uAjNzUjLwYzM14CM2MTN#top_anchor
func (f *Folder) rootFolder() *Folder {
	if f.Err != nil {
		return f
	}
	u := f.baseClient.urlJoin("/open-apis/drive/explorer/v2/root_folder/meta")
	_req, _ := http.NewRequest(http.MethodGet, u, nil)
	res := &RootFolderResp{}
	_, err := f.baseClient.CommonReq(_req, res)
	if err != nil {
		f.Err = newErr(err.Error())
		return f
	}
	return newFolder(res.Token, f.baseClient)
}

type RespCopyFile struct {
	FolderToken string `json:"folderToken"`
	Revision    int    `json:"revision"`
	Token       string `json:"token"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

type RespCreateFile = RespCreateFoler

type RespCreateFoler struct {
	URL      string `json:"url"`
	Revision int    `json:"revision"`
	Token    string `json:"token"`
}

type FolderMeta struct {
	ID        string `json:"id"` // id of the folder
	Name      string `json:"name"`
	Token     string `json:"token"`
	CreateUID string `json:"createUid"` // user id of the create user id
	EditUID   string `json:"editUid"`
	ParentID  string `json:"parentId"`
	OwnUID    string `json:"ownUid"`
}

type ChildrenInfo struct {
	ParentToken string              `json:"parentToken"`
	Children    map[string]Children `json:"children"`
}

type Children struct {
	Token string   `json:"token"`
	Name  string   `json:"name"`
	Type  FileType `json:"type"`
}

type (
	RootFolderResp struct {
		Token  string `json:"token"`
		ID     string `json:"id"`
		UserID string `json:"user_id"`
	}
)
