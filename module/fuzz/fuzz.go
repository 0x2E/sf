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

// FuzzModule 字典爆破模块
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
	logPrefix := "[" + f.Name + "]"

	// 加载字典
	dict, err := loadDict(app.Dict)
	if err != nil {
		return errors.Wrap(err, "failed to load dict file")
	}

	log.Printf("%s dict: %d\n", logPrefix, len(dict))

	// 设置泛解析黑名单
	err = f.Wildcard.Init(app)
	if err != nil {
		return errors.Wrap(err, "wildcard initialization failed")
	}
	log.Printf("%s wildcard initialization completed, blacklist: %d\n", logPrefix, len(f.Wildcard.Blacklist))

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
		return nil, errors.Wrap(err, "failed to open "+path)
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

	res := make([]string, len(dict))
	copy(res, dict) // 为了释放先前大cap的底层数组

	return res, nil
}
