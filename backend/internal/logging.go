package internal

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
}
