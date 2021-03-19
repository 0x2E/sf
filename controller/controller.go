package controller

import (
	"bufio"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/module"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

// Handle
func Handle(c *cli.Context) error {
	app := model.NewApp()

	if err := setup(app, c); err != nil { // 配置app
		return err
	}

	if err := module.Load(app); err != nil { // 运行模块
		return err
	}

	if err := output(app); err != nil { // 输出最终结果
		return err
	}

	log.Printf("Done, total time: %s\n", time.Since(app.Start))
	return nil
}

// setup 应用cli的参数
func setup(app *model.App, c *cli.Context) error {
	for _, f := range setAppList {
		err := f(app, c)
		if err != nil {
			return err
		}
	}
	return nil
}

// output 将结果输出到文件
func output(app *model.App) error {
	if len(app.Result.Data) == 0 { // 若没有结果则删除已创建的输出文件
		if err := os.Remove(app.Output); err != nil {
			log.Println("cannot delete empty output file: " + err.Error())
		}
		return nil
	}

	f, err := os.OpenFile(app.Output, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f) // 默认缓冲4096
	for k := range app.Result.Data {
		_, _ = bufWriter.WriteString(k + "\n")
	}
	_ = bufWriter.Flush()
	return nil
}
