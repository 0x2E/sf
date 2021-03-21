package model

import (
	"sync"
	"time"
)

// App 主体结构
type App struct {
	Domain   string    // 域名
	Dict     string    // 用于爆破的字典
	Resolver string    // DNS解析服务器
	Output   string    // 输出文件
	Thread   int       // 并发数
	Queue    int       // UDP发送-接收队列大小
	Wildcard int       // 泛解析处理模式，1 => 宽松模式，2 => 严格模式
	Start    time.Time // 开始时间
	Result   struct {  // 结果
		Data map[string]string
		Mu   sync.Mutex
	}
}

// NewApp 初始化一个新的App结构体
func NewApp() *App {
	return &App{
		Start: time.Now(),
		Result: struct {
			Data map[string]string
			Mu   sync.Mutex
		}{Data: make(map[string]string)},
	}
}
