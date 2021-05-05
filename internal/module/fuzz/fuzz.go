package fuzz

import (
	"bufio"
	"embed"
	"github.com/0x2E/sf/internal/module/fuzz/wildcard"
	"github.com/0x2E/sf/internal/option"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
	"sync"
)

// FuzzModule 字典爆破模块主体结构
type FuzzModule struct {
	name     string
	wildcard *wildcard.WildcardModel
	option   struct {
		domain   string
		dict     string
		resolver string
		thread   int
		queue    int
		retry    int
	}
	unReceived struct {
		data []string
		mu   sync.Mutex
	}
	result struct {
		data map[string]string
		mu   sync.Mutex
	}
}

// New 初始化一个fuzz模块
func New(option *option.Option) *FuzzModule {
	return &FuzzModule{
		name:     "fuzz",
		wildcard: wildcard.New(option.Wildcard.Mode, option.Wildcard.BlacklistMaxLen),
		option: struct {
			domain   string
			dict     string
			resolver string
			thread   int
			queue    int
			retry    int
		}{domain: option.Domain, dict: option.Dict, resolver: option.Resolver, thread: option.Thread, queue: option.Queue, retry: option.Retry},
		unReceived: struct {
			data []string
			mu   sync.Mutex
		}{data: make([]string, 0, 5000)},
		result: struct {
			data map[string]string
			mu   sync.Mutex
		}{data: make(map[string]string)},
	}
}

// GetName 返回名称
func (f *FuzzModule) GetName() string { return f.name }

// GetResult 返回结果
func (f *FuzzModule) GetResult() map[string]string { return f.result.data }

// Run 运行
func (f *FuzzModule) Run() error {
	logPrefix := "[" + f.name + "]"

	// 设置泛解析黑名单
	err := f.wildcard.Init(f.option.queue, f.option.domain, f.option.resolver)
	if err != nil {
		return errors.Wrap(err, "wildcard initialization failed")
	}
	log.Printf("%s wildcard blacklist: %d\n", logPrefix, f.wildcard.BlacklistLen())

	// 加载字典
	err = loadDict(f)
	if err != nil {
		return errors.Wrap(err, "failed to load dict file")
	}

	for try := 1; try <= (f.option.retry + 1); try++ {
		log.Printf("%s run #%d, queue remaining: %d\n", logPrefix, try, len(f.unReceived.data))

		ch := make(chan string, f.option.thread) // producer => consumer
		var wg sync.WaitGroup                    // producer(1) + consumer(n)

		wg.Add(1)
		go producer(ch, &wg, f)

		for i := 0; i < f.option.thread; i++ {
			wg.Add(1)
			go consumer(f, ch, &wg)
		}

		wg.Wait()

		if len(f.unReceived.data) == 0 {
			break
		}
	}

	return nil
}

type sfFile interface {
	Read(p []byte) (n int, err error)
	Close() error
}

//go:embed dict.txt
var embedFile embed.FS

// loadDict 加载字典
func loadDict(f *FuzzModule) error {
	var fs sfFile
	var err error
	if f.option.dict != "" {
		fs, err = os.Open(f.option.dict)
	} else {
		fs, err = embedFile.Open("dict.txt")
	}

	if err != nil {
		return errors.Wrap(err, "failed to open dict file")
	}
	defer fs.Close()

	// 字典去重
	suffix := "." + f.option.domain
	existMark := make(map[string]struct{}) // 标记已存在的数据
	scanner := bufio.NewScanner(fs)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		item := strings.TrimSpace(scanner.Text()) + suffix
		if _, ok := existMark[item]; ok {
			continue
		}
		f.unReceived.data = append(f.unReceived.data, item)
	}

	return nil
}
