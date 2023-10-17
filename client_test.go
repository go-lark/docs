package docs

import (
	"github.com/sirupsen/logrus"
)

func init() {
	SetLogLevel(logrus.DebugLevel)
}

func getClient() *Client {
	return NewClient(testAPPID, testAPPSecret, WithDomain(feishuDomain))
}
