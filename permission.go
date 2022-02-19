package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	memberTypeEmail            = "email"
	memberTypeOpenID           = "openid"
	memberTypeOpenChat         = "openchat"
	memberTypeUserID           = "userid"
	memberTypeOpendepartmentID = "opendepartmentid"
)

type permission struct {
	baseClient *Client
}

// Add for adding permission for a document
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uQDN5UjL0QTO14CN0kTN
// 增加权限
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uMzNzUjLzczM14yM3MTN
func (p *permission) Add(fileToken string, fileType FileType, perm Perm, notify bool, members ...*Member) error {
	for _, member := range members {
		reqBody := map[string]interface{}{
			"member_type": member.MemberType,
			"member_id":   member.MemberID,
			"perm":        perm,
		}
		path := fmt.Sprintf("/open-apis/drive/v1/permissions/%s/members?type=%s&need_notification=%t", fileToken, fileType, notify)
		url := p.baseClient.urlJoin(path)
		marshal, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshal))
		_, err := p.baseClient.CommonReq(req, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// TransferOwner for transferring the ownership.
// reference https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uUDN5UjL1QTO14SN0kTN
// 转移拥有者
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uQzNzUjL0czM14CN3MTN
func (p *permission) TransferOwner(fileToken string, fileType FileType, newOwner *Member, removeOldOwner, notify bool) ([]byte, *TransferOwnerResp, error) {
	reqBody := transferOwnerArgs{
		Type:           fileType,
		Token:          fileToken,
		Owner:          *newOwner,
		RemoveOldOwner: removeOldOwner,
		CancelNotify:   !notify,
	}
	switch newOwner.MemberType {
	case memberTypeEmail, memberTypeOpenID, memberTypeUserID:
	default:
		return nil, nil, newErr("owner user type is illegal, type: %s", newOwner.MemberType)
	}
	url := p.baseClient.urlJoin("/open-apis/drive/permission/member/transfer")
	marshal, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshal))
	resp := &TransferOwnerResp{}
	b, err := p.baseClient.CommonReq(req, resp)
	if err != nil {
		return nil, nil, err
	}
	return b, resp, nil
}

// PublicSet update document public setting
// 文档公共设置
// 参考 https://open.feishu.cn/document/ukTMukTMukTM/uITN5UjLyUTO14iM1kTN
func (p *permission) PublicSet(args *PublicSet) error {
	u := "/open-apis/drive/permission/v2/public/update/"
	u = p.baseClient.urlJoin(u)
	en, _ := json.Marshal(args)
	req, _ := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(en))
	_, err := p.baseClient.CommonReq(req, nil)
	return err
}

func NewMemberWithEmail(email string) *Member {
	return &Member{
		MemberType: memberTypeEmail,
		MemberID:   email,
	}
}

// NewMemberWithOpenID
// Parameter
//  openID: It is open user id. you can get it here(https://open.feishu.cn/document/home/user-identity-introduction/how-to-get)
func NewMemberWithOpenID(openID string) *Member {
	return &Member{
		MemberType: memberTypeOpenID,
		MemberID:   openID,
	}
}

func NewMemberWithOpenChatID(openChatID string) *Member {
	return &Member{
		MemberType: memberTypeOpenChat,
		MemberID:   openChatID,
	}
}

func NewMemberWithUserID(userID string) *Member {
	return &Member{
		MemberType: memberTypeUserID,
		MemberID:   userID,
	}
}

func NewMemberOpenDepartmentid(opendepartmentid string) *Member {
	return &Member{
		MemberType: memberTypeOpendepartmentID,
		MemberID:   opendepartmentid,
	}
}

type Member struct {
	MemberType string `json:"member_type"`
	MemberID   string `json:"member_id"`
}

type Perm string

const (
	PermView Perm = "view"
	PermEdit Perm = "edit"
	PermFull Perm = "full_access"
)

type (
	AddPermissionResp struct {
		IsAllSuccess bool                           `json:"is_all_success"`
		FailMembers  []AddPermissionRespFailMembers `json:"fail_members"`
	}

	AddPermissionRespFailMembers struct {
		MemberType string `json:"member_type"`
		MemberID   string `json:"member_id"`
		Perm       string `json:"perm"`
	}
)

type (
	transferOwnerArgs struct {
		Type           FileType `json:"type"`
		Token          string   `json:"token"`
		Owner          Member   `json:"owner"`
		RemoveOldOwner bool     `json:"remove_old_owner"`
		CancelNotify   bool     `json:"cancel_notify"`
	}

	TransferOwnerResp struct {
		IsSuccess bool                   `json:"is_success"`
		Token     string                 `json:"token"`
		Type      string                 `json:"type"`
		Owner     TransferOwnerRespOwner `json:"owner"`
	}

	TransferOwnerRespOwner struct {
		MemberType string `json:"member_type"`
		MemberID   string `json:"member_id"`
	}

	PublicSet struct {
		Token           string  `json:"token"`
		Type            string  `json:"type"`
		SecurityEntity  *string `json:"security_entity,omitempty"`
		CommentEntity   *string `json:"comment_entity,omitempty"`
		ShareEntity     *string `json:"share_entity,omitempty"`
		LinkShareEntity *string `json:"link_share_entity,omitempty"`
		ExternalAccess  *bool   `json:"external_access,omitempty"`
		InviteExternal  *bool   `json:"invite_external,omitempty"`
	}
)
