package fuzz

import (
	"bufio"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/module/fuzz/wildcard"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
	"sync"
)

// FuzzModule 字典爆破模块主体结构
type FuzzModule struct {
	Name       string
	Wildcard   *wildcard.WildcardModel
	UnReceived struct {
		Data []string
		Mu   sync.Mutex
	}
	Result struct {
		Data map[string]string
		Mu   sync.Mutex
	}
}

// New 初始化一个新的字典爆破模块结构体
func New() *FuzzModule {
	return &FuzzModule{
		Name:     "fuzz",
		Wildcard: &wildcard.WildcardModel{Blacklist: make(map[string]string)},
		UnReceived: struct {
			Data []string
			Mu   sync.Mutex
		}{Data: make([]string, 0, 5000)},
		Result: struct {
			Data map[string]string
			Mu   sync.Mutex
		}{Data: make(map[string]string)},
	}
}

// GetName 返回名称
func (f *FuzzModule) GetName() string { return f.Name }

// GetResult 返回结果
func (f *FuzzModule) GetResult() map[string]string { return f.Result.Data }

// Run 运行
func (f *FuzzModule) Run(app *model.App) error {
	logPrefix := "[" + f.Name + "]"

	// 设置泛解析黑名单
	if err := f.Wildcard.Init(app); err != nil {
		return errors.Wrap(err, "wildcard initialization failed")
	}
	log.Printf("%s wildcard initialization completed, blacklist: %d\n", logPrefix, len(f.Wildcard.Blacklist))

	// 加载字典
	if err := loadDict(app, f); err != nil {
		return errors.Wrap(err, "failed to load dict file")
	}

	for try := 1; try <= (app.Retry + 1); try++ {
		log.Printf("%s run#%d, dict: %d\n", logPrefix, try, len(f.UnReceived.Data))

		ch := make(chan string, app.Thread) // producer => consumer
		var wg sync.WaitGroup               // producer(1) + consumer(n)

		wg.Add(1)
		go producer(ch, &wg, f)

		for i := 0; i < app.Thread; i++ {
			wg.Add(1)
			go consumer(ch, &wg, app, f)
		}

		wg.Wait()

		if len(f.UnReceived.Data) == 0 {
			break
		}
	}

	return nil
}

// loadDict 加载字典
func loadDict(app *model.App, f *FuzzModule) error {
	fs, err := os.Open(app.Dict)
	if err != nil {
		return errors.Wrap(err, "failed to open dict file")
	}
	defer fs.Close()

	// 字典去重
	suffix := "." + app.Domain
	existMark := make(map[string]struct{}) // 标记已存在的数据
	scanner := bufio.NewScanner(fs)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		item := strings.TrimSpace(scanner.Text()) + suffix
		if _, ok := existMark[item]; ok {
			continue
		}
		f.UnReceived.Data = append(f.UnReceived.Data, item)
	}

	return nil
}
