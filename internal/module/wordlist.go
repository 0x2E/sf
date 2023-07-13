package module

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/0x2E/sf/internal/conf"
	"github.com/pkg/errors"
)

func RunWordlist(ctx context.Context, toNext chan<- *Task) error {
	if conf.C.Wordlist == "" {
		return errors.New("no wordlist input")
	}

	f, err := os.Open(conf.C.Wordlist)
	if err != nil {
		return errors.Wrap(err, "open wordlist file")
	}
	defer f.Close()

	var fn func(string) string
	if conf.C.RawTarget != conf.C.Target {
		fn = func(word string) string {
			return strings.ReplaceAll(conf.C.RawTarget, "%", word)
		}
	} else {
		suffix := "." + conf.C.Target
		fn = func(word string) string {
			return word + suffix
		}
	}

	dSet := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if _, ok := dSet[word]; ok {
			continue
		}
		dSet[word] = struct{}{}

		dn := fn(word)
		putTask(toNext, dn)
	}
	return nil
}
