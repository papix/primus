package server

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func init() {
	// init logger
	LogAccess = logrus.New()
	LogError = logrus.New()
}

func SetLogOut(log *logrus.Logger, outString string) error {
	switch outString {
	case "stdout":
		log.Out = os.Stdout
	case "stderr":
		log.Out = os.Stderr
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		log.Out = f
	}
	return nil
}

func SetLogLevel(log *logrus.Logger, levelString string) error {
	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		return err
	}
	log.Level = level
	return nil
}
