package utility

import (
	"io"
	"log"
	"os"
)

func fn1() {

}

func InitializeLogger() error {

	logdir := os.Getenv("LOG_DIR")
	if logdir == "" {
		logdir = "."
	}
	filename := logdir + string(os.PathSeparator) + "timetable_srv.log"
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return nil
}
