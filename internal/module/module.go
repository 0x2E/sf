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
	//todo 优先级字段，比如域传送的比爆破的优先级更高，域传送的结果应该覆盖爆破结果
}

// NewTask 创建一个任务并发送到目标队列
func NewTask(domain string, toNext chan<- *Task) {
	toNext <- &Task{
		Subdomain: dns.Fqdn(domain),
	}
}
