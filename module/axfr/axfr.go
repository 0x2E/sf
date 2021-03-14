package axfr

import (
	"github.com/0x2E/sf/model"
	"github.com/miekg/dns"
	"regexp"
	"strings"
	"sync"
)

// 测试环境：https://github.com/vulhub/vulhub/tree/master/dns/dns-zone-transfer

// AxfrModel 域传送检测模块
type AxfrModel struct {
	Name   string
	Result struct {
		Data map[string]string
		Mu   sync.Mutex
	}
}

func New() *AxfrModel {
	return &AxfrModel{
		Name: "zone-transfer",
		Result: struct {
			Data map[string]string
			Mu   sync.Mutex
		}{Data: make(map[string]string)},
	}
}

func (a *AxfrModel) GetName() string { return a.Name }

func (a *AxfrModel) GetResult() map[string]string { return a.Result.Data }

func (a *AxfrModel) Run(app *model.App) error {
	domain := dns.Fqdn(app.Domain)
	// 获取NS
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeNS)
	r, err := dns.Exchange(m, app.Resolver)
	if err != nil {
		return err
	}
	if len(r.Answer) == 0 {
		return nil
	}
	// 检测每个NS的域传送
	wg := sync.WaitGroup{}
	for _, v := range r.Answer {
		ns := strings.Replace(v.String(), v.Header().String(), "", 1)
		wg.Add(1)
		go transfer(a, &wg, domain, ns)
	}

	wg.Wait()
	return nil
}

func transfer(a *AxfrModel, wg *sync.WaitGroup, domain, ns string) {
	defer wg.Done()

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(domain)
	c, err := t.In(m, ns) // 默认2秒超时
	if err != nil {
		//fmt.Println(err.Error())
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

	a.Result.Mu.Lock()
	for k, v := range res {
		if _, ok := a.Result.Data[k]; !ok {
			a.Result.Data[k] = v
		}
	}
	a.Result.Mu.Unlock()
}
