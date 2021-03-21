package module

import (
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/module/axfr"
	"github.com/0x2E/sf/module/fuzz"
	"log"
	"sync"
	"time"
)

// moduleInterface 模块接口
type moduleInterface interface {
	Run(app *model.App) error     // 运行
	GetName() string              // 返回模块名称
	GetResult() map[string]string // 返回模块的结果
}

// Load 初始化、启动各个模块
func Load(app *model.App) {
	workers := []moduleInterface{
		moduleInterface(fuzz.New()),
		moduleInterface(axfr.New()),
	}
	var wg sync.WaitGroup // 各模块
	for i := range workers {
		wg.Add(1)
		go run(app, &wg, workers[i])
	}

	wg.Wait()
	//log.Println("all modules done")
}

// run 启动模块
func run(app *model.App, wg *sync.WaitGroup, worker moduleInterface) {
	defer wg.Done()

	startTime := time.Now()
	workerName := worker.GetName()
	logPrefix := "[" + workerName + "]"

	log.Println(logPrefix + " start")

	err := worker.Run(app)
	if err != nil {
		log.Printf("%s error: %s\n", logPrefix, err.Error())
		return
	}

	workerRes := worker.GetResult()
	log.Printf("%s done, subdomains: %d, time: %s\n", logPrefix, len(workerRes), time.Since(startTime))

	writeToApp(app, workerRes)
}

// writeToApp 将结果写入app.Result
func writeToApp(app *model.App, res map[string]string) {
	if len(res) == 0 {
		return
	}

	app.Result.Mu.Lock()
	for k, v := range res {
		if _, ok := app.Result.Data[k]; ok { // 结果已存在则跳过
			continue
		}
		app.Result.Data[k] = v
	}
	app.Result.Mu.Unlock()
	return
}
