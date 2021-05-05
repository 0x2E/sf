package module

import (
	"github.com/0x2E/sf/internal/module/axfr"
	"github.com/0x2E/sf/internal/module/fuzz"
	"github.com/0x2E/sf/internal/option"
)

// Module 模块接口
type Module interface {
	Run() error                   // 运行
	GetName() string              // 返回模块名称
	GetResult() map[string]string // 返回模块的结果
}

// Load 返回模块列表
func Load(option *option.Option) []Module {
	return []Module{
		Module(fuzz.New(option)),
		Module(axfr.New(option)),
	}
}
