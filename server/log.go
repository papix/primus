package server

import (
	"os"

	"github.com/Sirupsen/logrus"
	isatty "github.com/mattn/go-isatty"
)

func (ps *PrimusServer) SetupLogger() error {
	if err := ps.setupAccessLog(); err != nil {
		return err
	}
	if err := ps.setupErrorLog(); err != nil {
		return err
	}

	return nil
}

func (ps *PrimusServer) setupAccessLog() error {
	ps.AccessLog = logrus.New()
	filename := ps.Conf.Log.AccessLog

	if filename == "" {
		ps.AccessLog.Out = os.Stdout
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		ps.AccessLog.Out = f

		switch ps.AccessLog.Formatter.(type) {
		case *logrus.TextFormatter:
			ps.AccessLog.Formatter.(*logrus.TextFormatter).DisableColors = true
		}
	}

	return nil
}

func (ps *PrimusServer) setupErrorLog() error {
	ps.ErrorLog = logrus.New()
	filename := ps.Conf.Log.ErrorLog
	if filename == "" {
		ps.ErrorLog.Out = os.Stderr
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		ps.ErrorLog.Out = f

		switch ps.ErrorLog.Formatter.(type) {
		case *logrus.TextFormatter:
			ps.ErrorLog.Formatter.(*logrus.TextFormatter).DisableColors = true
		}
	}

	return nil
}

func (ps *PrimusServer) loggerClose() {
	ps.closeAccessLog()
	ps.closeErrorLog()
}

func (ps *PrimusServer) closeAccessLog() {
	switch ps.AccessLog.Out.(type) {
	case *os.File:
		if !isatty.IsTerminal(ps.AccessLog.Out.(*os.File).Fd()) {
			ps.AccessLog.Out.(*os.File).Close()
		}
	}
}

func (ps *PrimusServer) closeErrorLog() {
	switch ps.ErrorLog.Out.(type) {
	case *os.File:
		if !isatty.IsTerminal(ps.ErrorLog.Out.(*os.File).Fd()) {
			ps.ErrorLog.Out.(*os.File).Close()
		}
	}
}

// AccessLog family
func (ps *PrimusServer) Info(args ...interface{}) {
	ps.AccessLog.Info(args...)
}

func (ps *PrimusServer) Infoln(args ...interface{}) {
	ps.AccessLog.Infoln(args...)
}

func (ps *PrimusServer) Infof(format string, args ...interface{}) {
	ps.AccessLog.Infof(format, args...)
}

// ErrorLog family
func (ps *PrimusServer) Error(args ...interface{}) {
	ps.ErrorLog.Error(args...)
}

func (ps *PrimusServer) Errorln(args ...interface{}) {
	ps.ErrorLog.Errorln(args...)
}

func (ps *PrimusServer) Errorf(format string, args ...interface{}) {
	ps.ErrorLog.Errorf(format, args...)
}
