// docs 库封装了飞书开放平台的接口，提供便捷的飞书云文档操作能力。
package docs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-lark/docs/log"
	"github.com/hilaily/kit/httpx"
	"github.com/hilaily/kit/stringx"
	"github.com/sirupsen/logrus"
)

const (
	canUseTenantToken = 3 // if tenant token expired in 10 sec, use it, else renew it
	//canRenewTenantToken = 10 // if tenant token expired in 10 sec, renew it
)

// NewClient create a client with app id and app secret.
func NewClient(appID, appSecret string, ops ...ClientOption) *Client {
	c := &Client{
		appID:      appID,
		appSecret:  appSecret,
		domain:     "https://open.feishu.cn",
		httpClient: &http.Client{},
	}
	for _, f := range ops {
		f(c)
	}
	return c
}

// Client for docs sdk, support sheet, doc and so on
type Client struct {
	appID          string
	appSecret      string
	_token         string
	refreshToken   string
	tokenExpire    int64
	tenantAK       string
	tenantAKExpire int64
	domain         string
	httpClient     *http.Client
	useOldToken    bool

	tokenGetter func() (token string, expireTime int64, err error)
}

// CommonReq provide comment http request
func (c *Client) CommonReq(_req *http.Request, dst interface{}) ([]byte, error) {
	_req.Header.Set("Content-Type", "application/json")
	if log.IsDebugLevel() {
		var reqBody []byte
		reqURL := _req.URL.String()
		if _req.GetBody != nil {
			_body, err1 := _req.GetBody()
			if err1 != nil {
				log.Debugf("requeset, get body fail, %s", err1.Error())
			}
			reqBody, err1 = ioutil.ReadAll(_body)
			if err1 != nil {
				log.Debugf("common req,read body fail, %s", err1.Error())
			}
		}
		log.Debugf("request method: %s url: %s body: %s", _req.Method, reqURL, string(reqBody))
	}
	return c.DoRequest(_req, dst)
}

// DoRequest set authorization header and send request
func (c *Client) DoRequest(_req *http.Request, dst interface{}) ([]byte, error) {
	var token string
	var err error
	if c.tokenGetter != nil {
		token = c.getTokenFromTokenGetter()
	} else {
		token, err = c.getTenantAccessToken()
		if err != nil {
			return nil, fmt.Errorf("get tenant access token failed, %w", err)
		}
	}
	_req.Header.Set("Authorization", "Bearer "+token)
	reqURL := _req.URL.String()

	res, err := c.httpClient.Do(_req)
	if err != nil {
		return nil, fmt.Errorf("common req, url: %s, %w", reqURL, err)
	}
	defer res.Body.Close()
	body := &respBody{}
	b, err := httpx.HandleResp(res, body)
	if err != nil {
		return nil, fmt.Errorf("common req, handle resp, url: %s, %w", reqURL, err)
	}
	logrus.Debugf("response body: %s", string(b))
	if body.Code != 0 {
		return nil, &Err{Code: body.Code, Msg: body.Msg}
	}
	err = json.Unmarshal(body.Data, &dst)
	if err != nil {
		return nil, fmt.Errorf("unmarshal resp, url: %s, body: %s, %w", reqURL, string(body.Data), err)
	}
	return body.Data, nil
}

// SpreadSheet is for Sheets use
// Parameter
//  spreadSheetToken: token of a spreadsheets.
// Note
//  in a spreadsheet url, for example: https://abc.feishu.cn/sheets/shtcnjvusYPizPzZ8JqIWyCP7ca, shtcnjvusYPizPzZ8JqIWyCP7ca is the token
func (c *Client) OpenSpreadSheet(spreadSheetToken string) *SpreadSheet {
	ss := &SpreadSheet{}
	ss.baseClient = c
	ss.token = spreadSheetToken
	return ss
}

// Doc for doc operation
// Note
//  in a doc url, for example: https://abc.feishu.cn/docs/doccnuqdJJqnJ0LLWOjxoTS2Rld, doccnuqdJJqnJ0LLWOjxoTS2Rld is the token
func (c *Client) OpenDoc(token string) *Doc {
	d := &Doc{}
	d.baseClient = c
	d.token = token
	return d
}

// Folder for folder operation
// Note
//  in a folder url, for example: https://abc.feishu.cn/drive/folder/fldcnNhbqOyI0PVEPCuKa0acocdb, fldcnNhbqOyI0PVEPCuKa0acocdb is the token
func (c *Client) OpenFolder(token string) *Folder {
	f := &Folder{}
	f.baseClient = c
	f.token = token
	return f
}

// RootFolder get root folder of the bot/user
func (c *Client) RootFolder() *Folder {
	f := &Folder{}
	f.baseClient = c
	return f.rootFolder()
}

func (c *Client) file() *File {
	return &File{
		baseClient: c,
	}
}

func (c *Client) permission() *permission {
	return &permission{
		baseClient: c,
	}
}

func (c *Client) attachment() *Attachment {
	a := &Attachment{}
	a.f = c.file()
	return a
}

func (c *Client) getTokenFromTokenGetter() string {
	if c._token != "" && c.tokenExpire > time.Now().Unix() {
		return c._token
	} else if token, expire, err := c.tokenGetter(); err != nil {
		log.Errorf("get token from token getter err, %s\n", err.Error())
		return ""
	} else {
		c._token, c.tokenExpire = token, expire
		return c._token
	}
}

// ClientOption Client option parameters
type ClientOption func(c *Client)

// WithDomain set domain for api, default is https://open.feishu.cn
func WithDomain(domain string) ClientOption {
	return func(c *Client) {
		c.domain = strings.TrimSuffix(domain, "/")
	}
}

// WithProxy add http proxy
func WithProxy(proxyURL *url.URL, insecureSkipVerify bool) ClientOption {
	return func(c *Client) {
		c.httpClient.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		}
	}
}

// WithTimeout set http request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithTokenGetter set a function to getting the token.
func WithTokenGetter(f func() (token string, expireTime int64, err error)) ClientOption {
	return func(c *Client) {
		if f == nil {
			return
		}
		c.tokenGetter = f
	}
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func (c *Client) getTenantAccessToken() (string, error) {
	now := time.Now()
	liveSec := c.tenantAKExpire - now.Unix()
	if c.tenantAK != "" && liveSec > canUseTenantToken {
		//if liveSec < canRenewTenantToken {
		//	go func() {
		//		_ = c.getTenantAccessTokenBase()
		//	}()
		//}
		return c.tenantAK, nil
	}
	err := c.getTenantAccessTokenBase()
	return c.tenantAK, err
}
func (c *Client) getTenantAccessTokenBase() error {
	body, _ := json.Marshal(map[string]string{
		"app_id":     c.appID,
		"app_secret": c.appSecret,
	})
	_req, _ := http.NewRequest(http.MethodPost, c.urlJoin("/open-apis/auth/v3/tenant_access_token/internal/"), bytes.NewReader(body))
	respData := &respTenantAccessToken{}
	resp, err := c.httpClient.Do(_req)
	if err != nil {
		return fmt.Errorf("get tenant access token do http request failed, %w", err)
	}
	defer resp.Body.Close()
	_, err = httpx.HandleResp(resp, respData)
	if err != nil {
		return fmt.Errorf("get tenant access token failed, %w", err)
	}
	if respData.Code != 0 {
		return fmt.Errorf("get tenant access token failed, code: %d, msg: %s", respData.Code, respData.Msg)
	}
	c.tenantAK = respData.TenantAccessToken
	c.tenantAKExpire = respData.Expire
	return nil
}

func (c *Client) urlJoin(path ...string) string {
	return stringx.URLJoin(c.domain, path...)
}

type respTenantAccessToken struct {
	Code              int    `json:"code"`
	Expire            int64  `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}
