package oalog

import (
	"fmt"
	"github.com/tokopedia/tdk/v2/go/log"
	"gopkg.in/tokopedia/logging.v2"
	"os"
)

func New() error {
	pwd, _ := os.Getwd()
	errorFile := fmt.Sprintf("%s/log/error.log", pwd)
	debugFile := fmt.Sprintf("%s/log/access.log", pwd)
	cfg := &log.Config{
		Level:      "debug",
		LogFile:    errorFile,
		DebugFile:  debugFile,
		AppName:    "ordered-async-log",
		TimeFormat: "2006/01/02 15:04:05",
		Caller:     true,
		UseJSON:    true,
		CallerSkip: 3,
	}
	err := log.SetStdLog(cfg)
	if err != nil {
		return err
	}
	logging.LogInit()

	err = log.SetStdLog(cfg)
	if err != nil {
		return err
	}
	q = new(queue)
	q.msgChan = make(chan item)
	processQueue()
	return nil
}
