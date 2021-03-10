package fuzz

import (
	"bufio"
	"fmt"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/module/fuzz/wildcard"
	"github.com/0x2E/sf/util/logger"
	"os"
	"strings"
	"sync"
)

type FuzzModule struct {
	Name     string
	Wildcard *wildcard.WildcardModel
	Result   struct {
		Data map[string]string
		Mu   sync.Mutex
	}
}

func New() *FuzzModule {
	return &FuzzModule{
		Name:     "fuzz",
		Wildcard: &wildcard.WildcardModel{Blacklist: make(map[string]string)},
		Result: struct {
			Data map[string]string
			Mu   sync.Mutex
		}{Data: make(map[string]string)},
	}
}

func (f *FuzzModule) GetName() string { return f.Name }

func (f *FuzzModule) GetResult() map[string]string { return f.Result.Data }

func (f *FuzzModule) Run(app *model.App) error {
	// 加载字典
	dict, err := loadDict(app.Dict)
	if err != nil {
		return err
	}

	// 设置泛解析黑名单
	if err := f.Wildcard.Init(app); err != nil {
		return err
	}

	ch := make(chan string, app.Thread) // producer => consumer
	var wg sync.WaitGroup               // producer(1) + consumer(n)

	wg.Add(1)
	go producer(ch, &wg, dict)

	for i := 0; i < app.Thread; i++ {
		wg.Add(1)
		go consumer(ch, &wg, app, f)
	}

	wg.Wait()
	return nil
}

// loadDict 加载字典
func loadDict(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 字典去重
	dict := make([]string, 0, 50000)       // 去重后的字典，cap大一点减少底层数组扩容次数
	existMark := make(map[string]struct{}) // 标记已存在的数据
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		item := strings.TrimSpace(scanner.Text())
		if _, ok := existMark[item]; ok {
			continue
		}
		dict = append(dict, item)
	}
	res := make([]string, len(dict)) // 释放先前大cap的底层数组
	copy(res, dict)
	logger.Info(fmt.Sprintf("loaded entries from dict: %d", len(dict)))
	return res, nil
}
