package axfr

import (
	"github.com/0x2E/sf/internal/option"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"sync"
)

// 测试环境：https://github.com/vulhub/vulhub/tree/master/dns/dns-zone-transfer

// AxfrModel 域传送检测模块主体结构
type AxfrModel struct {
	name   string
	option struct {
		domain   string
		resolver string
	}
	result struct {
		data map[string]string
		mu   sync.Mutex
	}
}

// New 初始化域传送模块
func New(option *option.Option) *AxfrModel {
	return &AxfrModel{
		name: "zone-transfer",
		option: struct {
			domain   string
			resolver string
		}{domain: option.Domain, resolver: option.Resolver},
		result: struct {
			data map[string]string
			mu   sync.Mutex
		}{data: make(map[string]string)},
	}
}

// GetName 返回名称
func (a *AxfrModel) GetName() string { return a.name }

// GetResult 返回结果
func (a *AxfrModel) GetResult() map[string]string { return a.result.data }

// Run 运行
func (a *AxfrModel) Run() error {
	domain := dns.Fqdn(a.option.domain)

	// 获取NS
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeNS)
	r, err := dns.Exchange(m, a.option.resolver)
	if err != nil {
		return errors.Wrap(err, "failed to get NS record")
	}
	if len(r.Answer) == 0 {
		return nil
	}

	// 检测每个NS的域传送并保存结果
	wg := sync.WaitGroup{}
	for _, v := range r.Answer {
		ns := strings.Replace(v.String(), v.Header().String(), "", 1)
		wg.Add(1)
		go transfer(a, &wg, domain, ns)
	}

	wg.Wait()
	return nil
}

// transfer 检测传入的NS是否有域传送漏洞
func transfer(a *AxfrModel, wg *sync.WaitGroup, domain, ns string) {
	defer wg.Done()

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(domain)
	c, err := t.In(m, ns) // 默认2秒超时
	if err != nil {
		return
	}

	res := make(map[string]string) // 暂存结果
	// 匹配出域名和类型，若类型为SOA则不计入结果
	re, _ := regexp.Compile(`([\w.]*)\.\s*[0-9]+\s*(?:IN|CS|CH|HS)\s*(\w+)`)
	for v := range c { // 域传送用TCP，且可能不止一个包，每接收到一个都会写入channel
		if v.Error != nil {
			return
		}
		for _, rr := range v.RR {
			h := re.FindStringSubmatch(rr.Header().String())
			if len(h) != 3 {
				continue
			}
			if h[2] == "SOA" {
				continue
			}
			res[h[1]] = strings.Replace(rr.String(), rr.Header().String(), "", 1)
		}
	}

	a.result.mu.Lock()
	for k, v := range res {
		if _, ok := a.result.data[k]; !ok {
			a.result.data[k] = v
		}
	}
	a.result.mu.Unlock()
}
