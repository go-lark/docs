package log

import "github.com/sirupsen/logrus"

var (
	log = logrus.New()

	Errorf  = log.Errorf
	Errorln = log.Errorln
	Infof   = log.Infof
	Infoln  = log.Infoln
	Debugf  = log.Debugf
	Debugln = log.Debugln
)

func init() {
	log.SetLevel(logrus.FatalLevel)
}

func IsDebugLevel() bool {
	return log.Level == logrus.DebugLevel
}

func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}
