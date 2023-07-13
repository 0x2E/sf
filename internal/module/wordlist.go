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

	suffix := "." + conf.C.Target
	dSet := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		dn := strings.TrimSpace(scanner.Text()) + suffix
		if _, ok := dSet[dn]; ok {
			continue
		}
		dSet[dn] = struct{}{}
		putTask(toNext, dn)
	}
	return nil
}
