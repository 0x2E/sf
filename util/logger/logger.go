package logger

import (
	"log"
)

// Init 初始化logger
func Init() {
	log.SetFlags(log.Ltime)
}

// Info 输出信息
func Info(msg string) {
	log.Print(msg)
}

// Error 输出错误信息并终止程序
func Error(msg string) {
	log.Fatal("[Error]: " + msg)
}
