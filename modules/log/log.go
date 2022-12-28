package log

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

type ConfigOut interface {
}

func init() {
	logger = log.New(os.Stderr, "", log.LstdFlags)
}

func Println(v ...interface{}) {
	logger.Println(v)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v)

}

func Fatal(v ...interface{}) {
	logger.Fatal(v)

}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v)
}

func Panic(v ...interface{}) {
	logger.Panic(v)
}
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v)
}

func SaveLog(path string) {
	f, _ := os.OpenFile(path+"/run.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	multiWriter := io.MultiWriter(os.Stdout, f)
	logger.SetOutput(multiWriter)
}
