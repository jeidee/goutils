package log

import (
	"errors"
	"io"

	"os"

	log "github.com/Sirupsen/logrus"
)

// Log 로그 설정
type Log struct {
	Path      string `json:"path"`
	Level     string `json:"level"`
	Formatter string `json:"formatter"`
	Console   bool   `json:"console"`
}

type Logger struct {
	logger *log.Logger
}

var logger *Logger

func GetLogger(params ...interface{}) (*log.Logger, error) {
	if len(params) == 0 {
		if logger == nil {
			return nil, errors.New("Logger is not initialized")
		}
		return logger.logger, nil
	}

	log.SetFormatter(&log.JSONFormatter{})

	logger = &Logger{}
	logger.logger = log.StandardLogger()

	logConf := params[0].(*Log)

	f, err := os.OpenFile(logConf.Path, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Error(err)
	}
	if logConf.Console {
		log.SetOutput(io.MultiWriter(f, os.Stdout))
	} else {
		log.SetOutput(f)
	}

	level, err := log.ParseLevel(logConf.Level)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	if logConf.Formatter == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	return logger.logger, nil
}

func Debug(args ...interface{}) {
	logger.logger.Debug(args)
}

func Info(args ...interface{}) {
	logger.logger.Info(args)
}

func Warn(args ...interface{}) {
	logger.logger.Warn(args)
}

func Error(args ...interface{}) {
	logger.logger.Error(args)
}

func Fatal(args ...interface{}) {
	logger.logger.Fatal(args)
}

func Panic(args ...interface{}) {
	logger.logger.Panic(args)
}

func Debugf(format string, args ...interface{}) {
	logger.logger.Debugf(format, args)
}

func Infof(format string, args ...interface{}) {
	logger.logger.Infof(format, args)
}

func Warnf(format string, args ...interface{}) {
	logger.logger.Warnf(format, args)
}

func Errorf(format string, args ...interface{}) {
	logger.logger.Errorf(format, args)
}

func Fatalf(format string, args ...interface{}) {
	logger.logger.Fatalf(format, args)
}

func Panicf(format string, args ...interface{}) {
	logger.logger.Panicf(format, args)
}
