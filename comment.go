package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddWholeComment represent add a whole _comment
func (c comment) AddWholeComment(fileType FileType, content string) (*RespComment, error) {
	en, _ := json.Marshal(map[string]interface{}{
		"type":    fileType,
		"token":   c.fileToken,
		"content": content,
	})
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/open-apis/comment/add_whole", c.baseClient.domain),
		bytes.NewReader(en),
	)
	resp := &RespComment{}
	_, err := c.baseClient.CommonReq(req, resp)
	return resp, err
}

func newComment(token string, originClient *Client) *comment {
	return &comment{
		fileToken:  token,
		baseClient: originClient,
	}
}

// comment represent add comment for a doc
type comment struct {
	fileToken  string
	baseClient *Client
}

type RespComment struct {
	CommentID       string `json:"comment_id"`
	ReplyID         string `json:"reply_id"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}
