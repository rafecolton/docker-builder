package kamino

import "github.com/Sirupsen/logrus"

// Logger is the logger used by the runner package.  It is initialized in the
// init() function so it may be overwritten any time after that.
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.TextFormatter{}
	Logger.Level = logrus.InfoLevel // default to Info
}
