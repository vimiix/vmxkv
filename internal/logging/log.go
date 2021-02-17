package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	logger             = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logPrefix          string
	defaultCallerDepth = 2
	timeFormat         = "20060102"
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func Init(dir string) {
	fmt.Println("logging dir: " + dir)
	f := openLog(dir)
	logger = log.New(f, logPrefix, log.LstdFlags|log.Lshortfile)
}

func GetLogger() *log.Logger {
	return logger
}

func openLog(logDir string) (f *os.File) {
	_, err := os.Stat(logDir)
	switch {
	case os.IsNotExist(err):
		if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Fail to Mkdir : %v\n", err)
		}
	case os.IsPermission(err):
		log.Fatalf("Permission :%v\n", err)
	}

	fileName := fmt.Sprintf("vmxkv.%s.log", time.Now().Format(timeFormat))
	filePath := path.Join(logDir, fileName)

	f, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v\n", err)
	}

	return
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Output(defaultCallerDepth, fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Output(defaultCallerDepth, fmt.Sprintln(v...))
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Output(defaultCallerDepth, fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Output(defaultCallerDepth, fmt.Sprintln(v...))
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Output(defaultCallerDepth, fmt.Sprintln(v...))
	os.Exit(1)
}

func setPrefix(level Level) {
	logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])
	logger.SetPrefix(logPrefix)
}
