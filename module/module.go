package module

import (
	"fmt"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/module/axfr"
	"github.com/0x2E/sf/module/fuzz"
	"github.com/0x2E/sf/util/logger"
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
func Load(app *model.App) error {
	workers := []moduleInterface{
		moduleInterface(fuzz.New()),
		moduleInterface(axfr.New()),
	}
	var wg sync.WaitGroup // 各模块
	wg.Add(len(workers))
	for i, _ := range workers {
		go run(app, &wg, workers[i])
	}
	wg.Wait()
	logger.Info("all modules done")
	return nil
}

// run 启动模块
func run(app *model.App, wg *sync.WaitGroup, worker moduleInterface) {
	defer wg.Done()

	startTime := time.Now()
	workerName := worker.GetName()
	loggerPrefix := "module [" + workerName + "] "
	logger.Info(loggerPrefix + "start")

	err := worker.Run(app)
	if err != nil {
		logger.Error(loggerPrefix + "error and stop: " + err.Error())
		return
	}

	workerRes := worker.GetResult()
	logger.Info(fmt.Sprintf("%s done, subdomains: %d, time: %s", loggerPrefix, len(workerRes), time.Since(startTime)))
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
