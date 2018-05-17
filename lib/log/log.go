package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	lock   sync.RWMutex
	logger *log.Logger
	lgFile *os.File
	lgDir  string
	lgDay  string
)

func SetDir(logDir string) {
	lgDir = logDir
	err := os.MkdirAll(lgDir, 0755)
	if err != nil {
		panic(err)
	}
}

func Info(format string, v ...interface{}) {
	F("INFO "+format, v...)
}

func Warn(format string, v ...interface{}) {
	F("WARN "+format, v...)
}

func Err(format string, v ...interface{}) {
	F("ERROR "+format, v...)
}

func F(format string, v ...interface{}) {
	lock.Lock()
	checkLogFile()
	logger.Printf(format, v...)
	lock.Unlock()
}

func checkLogFile() {
	nowDay := time.Now().In(time.UTC).Format("2006-01-02")
	if lgDay == nowDay {
		return
	}

	lgDay = nowDay
	if lgFile != nil {
		lgFile.Close()
		fmt.Println("log file closed.")
	}
	lgFile, err := os.OpenFile(fmt.Sprintf("%s/%s.log", lgDir, lgDay), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open log file failed: %v\n", err)
	}
	logger = log.New(lgFile, "", log.LstdFlags|log.LUTC)
}
