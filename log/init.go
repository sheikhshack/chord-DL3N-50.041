package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func init() {
	_ = os.Mkdir("tmp", os.ModeDir+0777)
	logFileName := fmt.Sprintf("./tmp/log_%.23s.txt", time.Now().UTC())

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Info = log.New(file, "\u001b[34mINFO :\u001b[0m ", log.Ltime|log.Lshortfile)
	Warn = log.New(file, "\u001b[33mWARN :\u001b[0m ", log.Ltime|log.Lshortfile)
	Error = log.New(file, "\u001b[31mERROR:\u001b[0m ", log.Ltime|log.Lshortfile)

	LOG_LEVEL := os.Getenv("LOG")
	if LOG_LEVEL == "" || LOG_LEVEL == "quiet" || LOG_LEVEL == "error" {
		LOG_LEVEL = "error"

		mw := io.MultiWriter(os.Stdout, file)
		Error.SetOutput(mw)
	} else if LOG_LEVEL == "warn" {
		mw := io.MultiWriter(os.Stdout, file)
		Warn.SetOutput(mw)
		Error.SetOutput(mw)
	} else if LOG_LEVEL == "info" {
		mw := io.MultiWriter(os.Stdout, file)
		Info.SetOutput(mw)
		Warn.SetOutput(mw)
		Error.SetOutput(mw)
	}
}

func PrintTest() {
	Info.Printf("this is info\n")
	Warn.Printf("this is warn\n")
	Error.Printf("this is error\n")
}
