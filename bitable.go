package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func newBitable(token string, client *Client) *Bitable {
	return &Bitable{
		baseClient: client,
		tokenIns: tokenIns{
			token: token,
		},
	}
}

// Doc represent a doc file
type Bitable struct {
	Err error
	tokenIns
	baseClient *Client
	_comment   *comment
}

func (b *Bitable) Table(id string) *Table {
	return &Table{
		id:      id,
		Bitable: b,
	}
}

type Table struct {
	Err error
	id  string
	*Bitable
}

func (t *Table) AddRecord(data []Field) *Table {
	_url := t.baseClient.urlJoin(fmt.Sprintf("/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_create", t.token, t.id))
	r := []*Record{}
	for _, v := range data {
		r = append(r, &Record{Fields: v})
	}
	en, _ := json.Marshal(
		map[string]interface{}{
			"records": r,
		})
	_req, _ := http.NewRequest(http.MethodPost, _url, bytes.NewReader(en))
	resp := &AddRecordResp{}
	_, err := t.Bitable.baseClient.CommonReq(_req, resp)
	if err != nil {
		t.Err = err
	}
	return t
}

type Field = map[string]interface{}

type Record struct {
	Fields   Field  `json:"fields"`
	ID       string `json:"id,omitempty"`
	RecordID string `json:"record_id,omitempty"`
}

type AddRecordResp struct {
	Records []Record `json:"records"`
}
