package module

import (
	"github.com/0x2E/sf/internal/conf"
	"github.com/miekg/dns"
)

// Module 模块接口
type Module interface {
	Run() error
	GetName() string
}

type base struct {
	name   string
	conf   *conf.Config
	toNext chan<- *Task
}

func (b *base) GetName() string { return b.name }

// Load 初始化所有模块，返回模块列表
func Load(conf *conf.Config, toEnumerator, toRecorder chan<- *Task) []Module {
	return []Module{
		newAxfr(conf, toRecorder),
		newWordlist(conf, toEnumerator),
	}
}

// Task 任务
type Task struct {
	Subdomain string
	Answer    []dns.RR
	Time      int64 // 发出DNS请求的时间
	Received  bool
	Valid     bool // 表示分别存入valid和invalid结果集。传入的checker的可能被修改，否则保持NewTask时的赋值直到recorder
	//todo 优先级字段，比如域传送的比爆破的优先级更高，域传送的结果应该覆盖爆破结果
}

// NewTask 创建一个任务并发送到目标队列
func NewTask(toNext chan<- *Task, domain string) {
	toNext <- &Task{
		Subdomain: dns.Fqdn(domain),
		Valid:     true,
	}
}
