package conf

import (
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"os"
)

const (
	TESTDOMAIN = "github.com"
)

type Config struct {
	Domain   string // 目标域名
	Wordlist string // 字典路径
	Resolver string // DNS服务器
	Thread   int    // enumerator并发数
	Rate     int    // 每秒最大发包数
	Retry    int    // 重试次数
	Check    bool   // 是否开启有效性检查
}

// Verify 检查参数是否正常
func (c *Config) Verify() error {
	if _, ok := dns.IsDomainName(c.Domain); !ok {
		return errors.New("invalid domain name")
	}
	c.Domain = dns.Fqdn(c.Domain)

	if c.Wordlist != "" {
		f, err := os.Open(c.Wordlist)
		if err != nil {
			return errors.Wrap(err, "open wordlist file")
		}
		f.Close()
	}

	c.Resolver = c.Resolver + ":53"
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(TESTDOMAIN), dns.TypeA)
	if _, err := dns.Exchange(m, c.Resolver); err != nil {
		return errors.Wrap(err, "resolver may be invalid")
	}

	if c.Thread < 0 || c.Retry < 0 {
		return errors.New("numerical parameters must be positive")
	}

	return nil
}
