package docs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-lark/docs/doctypes"
)

func newDoc(token string, client *Client) *Doc {
	return &Doc{
		baseClient: client,
		tokenIns: tokenIns{
			token: token,
		},
	}
}

// Doc represent a doc file
type Doc struct {
	Err error
	tokenIns
	baseClient *Client
	_comment   *comment
}

// GetMeta for doc meta info
func (d *Doc) GetMeta() (*DocMeta, error) {
	u := fmt.Sprintf("%s/open-apis/doc/v2/meta/%s", d.baseClient.domain, d.token)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	res := &DocMeta{}
	_, err := d.baseClient.CommonReq(req, res)
	return res, err
}

func (d *Doc) GetContent() ([]byte, *DocContent, int, error) {
	u := fmt.Sprintf("%s/open-apis/doc/v2/%s/content", d.baseClient.domain, d.token)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp := &respGetContent{}
	r, e := d.baseClient.CommonReq(req, resp)
	if e != nil {
		return nil, nil, 0, e
	}
	content := &DocContent{}
	err := json.Unmarshal([]byte(resp.Content), content)
	return r, content, resp.Revision, err
}

func (d *Doc) AddWholeComment(content string) (*RespComment, error) {
	return d.getComment().AddWholeComment(FileTypeDoc, content)
}

// Share to other user or group.
func (d *Doc) Share(perm Perm, notify bool, members ...*Member) *Doc {
	if d.Err != nil {
		return d
	}
	err := d.baseClient.permission().Add(d.token, FileTypeDoc, perm, notify, members...)
	d.Err = err
	return d
}

func (d *Doc) ChangeOwner(newOwner *Member, removeOldOwner, notify bool) *Doc {
	if d.Err != nil {
		return d
	}
	_, _, err := d.baseClient.permission().TransferOwner(d.token, FileTypeDoc, newOwner, removeOldOwner, notify)
	if err != nil {
		d.Err = err
	}
	return d
}

func (d *Doc) SetAccessPermission(per string) *Doc {
	return d.setPublic(&PublicSet{
		LinkShareEntity: &per,
	})
}

func (d *Doc) Statistics() (*FileStatistics, error) {
	return newFile(d.baseClient).statistics(d.token, FileTypeDoc)
}

func (d *Doc) setPublic(args *PublicSet) *Doc {
	if d.Err != nil {
		return d
	}
	err := d.baseClient.permission().PublicSet(args)
	d.Err = err
	return d
}

type DocMeta struct {
	CreateDate     string        `json:"create_date"`
	CreateTime     int           `json:"create_time"`
	CreateUID      string        `json:"create_uid"`
	CreateUserName string        `json:"create_user_name"`
	DeleteFlag     DocDeleteFlag `json:"delete_flag"`
	EditTime       int           `json:"edit_time"`
	EditUserName   string        `json:"edit_user_name"`
	IsExternal     bool          `json:"is_external"`
	IsPined        bool          `json:"is_pined"`
	IsStared       bool          `json:"is_stared"`
	ObjType        string        `json:"obj_type"`
	OwnerID        string        `json:"owner_id"`
	OwnerUserName  string        `json:"owner_user_name"`
	ServerTime     int           `json:"server_time"`
	TenantID       string        `json:"tenant_id"`
	Title          string        `json:"title"`
	Type           int           `json:"type"`
	URL            string        `json:"url"`
}

type DocDeleteFlag int

var (
	DocDeleteFlagNormal  DocDeleteFlag = 0
	DocDeleteFlagTrashed DocDeleteFlag = 1
	DocDeleteFlagDeleted DocDeleteFlag = 2
)

func (d DocDeleteFlag) MarshalJSON() ([]byte, error) {
	return []byte{byte(int(d) - 1 + int('1'))}, nil
}

func (d *DocDeleteFlag) UnmarshalJSON(data []byte) error {
	*d = DocDeleteFlag(data[0] + 1 - '1')
	return nil
}

type DocContent struct {
	Title *doctypes.Title `json:"title"`
	Body  *doctypes.Body  `json:"body"`
}

type RespCreateDoc struct {
	ObjToken string `json:"objToken"`
	URL      string `json:"url"`
}

type respGetContent struct {
	Content  string `json:"content"`
	Revision int    `json:"revision"`
}

func (d *Doc) getComment() *comment {
	if d._comment == nil {
		d._comment = newComment(d.token, d.baseClient)
	}
	return d._comment
}
