package log

import "github.com/sirupsen/logrus"

var (
	Logger = logrus.New()

	Errorf  = Logger.Errorf
	Errorln = Logger.Errorln
	Infof   = Logger.Infof
	Infoln  = Logger.Infoln
	Debugf  = Logger.Debugf
	Debugln = Logger.Debugln
)

func init() {
	Logger.SetLevel(logrus.FatalLevel)
}

func IsDebugLevel() bool {
	return Logger.Level == logrus.DebugLevel
}

func SetLevel(level logrus.Level) {
	Logger.SetLevel(level)
}
