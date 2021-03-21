package controller

import (
	"bufio"
	"fmt"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/module"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

// handle 控制整个app的运行周期
func handle(c *cli.Context) error {
	app := model.NewApp()

	setup(app, c) // 配置app

	module.Load(app) // 运行模块

	output(app) // 输出最终结果

	log.Printf("done, total subdomains: %d, total time: %s\n", len(app.Result.Data), time.Since(app.Start))
	return nil
}

// setup 应用cli的参数
func setup(app *model.App, c *cli.Context) {
	for _, f := range setAppList {
		f(app, c)
	}
}

// output 将结果输出到文件
func output(app *model.App) {
	if len(app.Result.Data) == 0 {
		return
	}

	f, err := os.OpenFile(app.Output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil { // 无法创建结果文件，就输出到终端
		log.Print("failed to write results into file, so output to the console: " + err.Error())
		fmt.Println("============")
		for k := range app.Result.Data {
			fmt.Println(k)
		}
		return
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f) // 默认缓冲4096
	for k := range app.Result.Data {
		_, _ = bufWriter.WriteString(k + "\n")
	}
	_ = bufWriter.Flush()
}
