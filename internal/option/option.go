package option

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"os"
	"regexp"
)

type Option struct {
	Domain   string // 域名
	Dict     string // 用于fuzz的字典
	Resolver string // DNS解析服务器
	Output   string // 输出文件
	Thread   int    // fuzz并发数
	Queue    int    // fuzz任务队列长度
	Retry    int    // 重试次数
	Wildcard struct {
		Mode            int // 处理模式，1 => 宽松模式，2 => 严格模式
		BlacklistMaxLen int // 黑名单最大程度
	} // 泛解析
}

// ParseOption 解析并校验命令行参数
func ParseOption() (*Option, error) {
	o := &Option{}
	flag.StringVar(&o.Domain, "u", "", "Target url or domain name")
	flag.StringVar(&o.Dict, "d", "", "Load dictionary from a file")
	flag.StringVar(&o.Output, "o", "", "Output results to a file")
	flag.StringVar(&o.Resolver, "r", "8.8.8.8", "[fuzz] DNS resolver")
	flag.IntVar(&o.Thread, "t", 100, "[fuzz] The number of threads. Each thread will occupy a temporary port of the system until the end of the fuzz")
	flag.IntVar(&o.Queue, "q", 100, "[fuzz] The length of the task queue. Too high may fill the system socket buffer and cause packet loss")
	flag.IntVar(&o.Wildcard.Mode, "w", 1, "[fuzz] Two modes (1 or 2) for processing wildcard records. Mode 1 is only based on the IP blacklist. Mode 2 matches the IP blacklist, compares the similarity of web page titles after hits, and degenerates to mode 1 if port 80 cannot be accessed")
	flag.IntVar(&o.Wildcard.BlacklistMaxLen, "wl", 1000, "[fuzz] The maximum length of the IP blacklist for wildcard records")
	flag.IntVar(&o.Retry, "R", 2, "[fuzz] The number of retries")

	flag.Parse()

	if o.Domain == "" {
		return nil, errors.New("domain name cannot be empty")
	} else {
		re := regexp.MustCompile(`^(?:\w*://)?((?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6})/?`)
		domain := re.FindStringSubmatch(o.Domain)
		if len(domain) == 1 || len(domain[1]) < 4 {
			return nil, errors.New("invalid domain name")
		}
		o.Domain = domain[1]
	}

	if o.Dict != "" {
		f, err := os.Open(o.Dict)
		if err != nil {
			return nil, err
		}
		f.Close()
	}

	if o.Output == "" {
		o.Output = o.Domain + ".txt"
	}

	o.Resolver = o.Resolver + ":53"
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	if _, err := dns.Exchange(m, o.Resolver); err != nil {
		return nil, errors.New("resolver may be invalid: " + err.Error())
	}

	if o.Wildcard.Mode != 1 && o.Wildcard.Mode != 2 {
		return nil, errors.New("wildcard mode must be 1 or 2")
	}

	if o.Thread < 0 || o.Queue < 0 || o.Wildcard.BlacklistMaxLen < 0 || o.Retry < 0 {
		return nil, errors.New("numerical parameters cannot be negative")
	}

	return o, nil
}
