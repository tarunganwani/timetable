package util

import (
	"io"
	"log"
	"os"
)

func InitializeLogger(logname string) error {

	logdir := os.Getenv("LOG_DIR")
	if logdir == "" {
		logdir = "."
	}
	filename := logdir + string(os.PathSeparator) + logname
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return nil
}
