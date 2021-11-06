package module

import (
	"bufio"
	"github.com/0x2E/sf/internal/conf"
	"github.com/pkg/errors"
	"os"
	"strings"
)

// Wordlist 字典模块
type Wordlist struct {
	base
}

func newWordlist(conf *conf.Config, toEnumerator chan<- *Task) *Wordlist {
	return &Wordlist{
		base: base{
			name:   "wordlist",
			conf:   conf,
			toNext: toEnumerator,
		},
	}
}

func (d *Wordlist) Run() error {
	if d.conf.Wordlist == "" {
		return errors.New("no wordlist input")
	}
	f, err := os.Open(d.conf.Wordlist)

	if err != nil {
		return errors.Wrap(err, "open wordlist file")
	}
	defer f.Close()

	// 字典去重
	suffix := "." + d.conf.Domain
	existMark := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		item := strings.TrimSpace(scanner.Text()) + suffix
		if _, ok := existMark[item]; ok {
			continue
		}
		NewTask(d.toNext, item)
	}
	return nil
}
