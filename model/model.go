package model

import (
	"sync"
	"time"
)

// App 主体结构
type App struct {
	Domain   string // 域名
	Dict     string // 用于fuzz的字典
	Resolver string // DNS解析服务器
	Output   string // 输出文件
	Thread   int    // fuzz并发数
	Queue    int    // fuzz任务队列长度
	Retry    int    // 重试次数
	Wildcard struct {
		Mode            int // 处理模式，1 => 宽松模式，2 => 严格模式
		BlacklistMaxLen int // 黑名单最大程度
	} // 泛解析

	Start  time.Time // 开始时间
	Result struct {  // 汇总的结果
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
