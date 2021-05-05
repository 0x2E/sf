package engine

import (
	"bufio"
	"fmt"
	"github.com/0x2E/sf/internal/module"
	"github.com/0x2E/sf/internal/option"
	"log"
	"os"
	"sync"
	"time"
)

type Engine struct {
	Option *option.Option
	Start  time.Time
	Result struct {
		Data map[string]string
		Mu   sync.Mutex
	}
}

// New 初始化引擎
func New(o *option.Option) *Engine {
	return &Engine{
		Option: o,
		Start:  time.Now(),
		Result: struct {
			Data map[string]string
			Mu   sync.Mutex
		}{Data: make(map[string]string)},
	}
}

// Run 运行
func (e *Engine) Run() error {
	wg := sync.WaitGroup{}
	for _, i := range module.Load(e.Option) {
		wg.Add(1)
		go func(i module.Module, wg *sync.WaitGroup) {
			defer wg.Done()

			startTime := time.Now()
			logPrefix := "[" + i.GetName() + "]"

			log.Println(logPrefix + " start")

			err := i.Run()
			if err != nil {
				log.Printf("%s error: %s\n", logPrefix, err.Error())
				return
			}

			// 保存本模块的结果
			res := i.GetResult()
			log.Printf("%s done, subdomains: %d, time: %s\n", logPrefix, len(res), time.Since(startTime))
			if len(res) != 0 {
				e.Result.Mu.Lock()
				for k, v := range res {
					if _, ok := e.Result.Data[k]; ok { // 结果已存在则跳过
						continue
					}
					e.Result.Data[k] = v
				}
				e.Result.Mu.Unlock()
			}
		}(i, &wg)
	}
	wg.Wait()

	e.output()
	return nil
}

// output 输出最终结果
func (e *Engine) output() {
	if len(e.Result.Data) == 0 {
		return
	}

	f, err := os.OpenFile(e.Option.Output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil { // 无法创建结果文件，就输出到终端
		log.Print("failed to write results into file, so output to the console: " + err.Error())
		fmt.Println("============")
		for k := range e.Result.Data {
			fmt.Println(k)
		}
		return
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f) // 默认缓冲4096
	for k := range e.Result.Data {
		_, _ = bufWriter.WriteString(k + "\n")
	}
	_ = bufWriter.Flush()
}
