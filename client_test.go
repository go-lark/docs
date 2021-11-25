package docs

import (
	"github.com/sirupsen/logrus"
)

func init() {
	SetLogLevel(logrus.DebugLevel)
}

// Deprecated
//func getClient() *Client {
//	return NewClientWithRefreshToken(testAPPID, testAPPSecret, "")
//}

func getClientNew() *Client {
	return NewClient(testAPPID2, testAPPSecret2)
}

//func TestRefreshToken(t *testing.T) {
//	c := getClient()
//	refrshToken := "ur-pW6g5NuMu2HaGmr5Dfh4Mf"
//	err := c.GetTokenByRefresh(refrshToken)
//	assert.NoError(t, err)
//}
//
//func TestGetTokenByCode(t *testing.T) {
//	appID := os.Getenv("TEST_BOT_APP_ID")
//	appSecret := os.Getenv("TEST_BOT_APP_SECRET")
//	code := ""
//	c := NewClientWithRefreshToken(appID, appSecret, code)
//	res, err := c.GetTokenByCode(appID, appSecret, code)
//	assert.NoError(t, err)
//	t.Log(res)
//}

//func TestRefreshToken(t *testing.T) {
//	appID := os.Getenv("TEST_BOT_APP_ID")
//	appSecret := os.Getenv("TEST_BOT_APP_SECRET")
//	token := ""
//	c := NewClient(appID, appSecret, token)
//	err := c.getTokenByRefresh(c.refreshToken)
//	assert.NoError(t, err)
//	spew.Dump(c)
//}
