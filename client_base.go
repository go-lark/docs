package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hilaily/kit/httpx"
)

func (c *Client) GetMeta(token, typ, userIDType string) (*MetaRespMetas, error) {
	u := fmt.Sprintf("%s/open-apis/drive/v1/metas/batch_query", c.domain)
	if userIDType != "" {
		httpx.AddParamsToURL(u, &url.Values{"user_id_type": []string{userIDType}})
	}
	en, _ := json.Marshal(map[string]interface{}{
		"request_docs": []map[string]string{{
			"doc_token": token,
			"doc_type":  typ,
		}},
		"with_url": true,
	})
	req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(en))
	res := &MetaResp{}
	_, err := c.CommonReq(req, res)
	if err != nil {
		return nil, err
	}
	if len(res.FailedList) > 0 {
		return nil, fmt.Errorf("get meta fail, code:%d", res.FailedList[0].Code)
	}
	return res.Metas[0], nil
}

type MetaResp struct {
	Metas      []*MetaRespMetas      `json:"metas"`
	FailedList []*MetaRespFailedList `json:"failed_list"`
}

type MetaRespMetas struct {
	DocToken         string `json:"doc_token"`
	DocType          string `json:"doc_type"`
	Title            string `json:"title"`
	OwnerID          string `json:"owner_id"`
	CreateTime       int64  `json:"create_time,string"`
	LatestModifyUser string `json:"latest_modify_user"`
	LatestModifyTime int64  `json:"latest_modify_time,string"`
	URL              string `json:"url"`
	SecLabelName     string `json:"sec_label_name"`
}

type MetaRespFailedList struct {
	Token string `json:"token"`
	Code  int    `json:"code"`
}
