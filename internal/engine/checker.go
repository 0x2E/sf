package engine

import (
	"github.com/miekg/dns"
	"sync"
)

type check struct {
	existWildcard bool
	//existHTTP     bool
	//httpMark      []*req
	dnsMark []dns.RR
}

// existWildcard 判断是否存在泛解析
func (e *Engine) existWildcard() bool {
	m := new(dns.Msg)
	m.SetQuestion(e.conf.Domain, dns.TypeNS)
	r, err := dns.Exchange(m, e.conf.Resolver)
	if err != nil || r.Rcode != dns.RcodeSuccess || len(r.Answer) == 0 {
		return false
	}
	var n *dns.NS
	var ok bool
	for _, v := range r.Answer {
		if n, ok = v.(*dns.NS); !ok {
			continue
		}
		m := &dns.Msg{}
		m.SetQuestion("*."+e.conf.Domain, dns.TypeA)
		wResp, err := dns.Exchange(m, n.Ns+":53")
		if err != nil || wResp.Rcode != dns.RcodeSuccess || len(wResp.Answer) == 0 {
			continue
		}
		e.check.dnsMark = wResp.Answer
		break
	}
	if len(e.check.dnsMark) == 0 {
		return false
	}
	e.check.dnsMark[0].Header().Name = "" // 将DNS查询的域名置空，方便checker比较

	//e.check.httpMark, err = getHTTP(fmt.Sprintf("http://%s.%s", util.RandString(8), e.conf.Domain))
	//if err == nil && len(e.check.httpMark) != 0 {
	//	e.check.existHTTP = true
	//}
	return true
}

// checker 负责判定子域名是否有效，目前用于筛选掉泛解析域名。原本实现了两个判别条件（DNS特征和HTTP特征），但HTTP的性价比实在太低，暂时注释掉。
//
// 考虑到负载均衡、网关架构等情况的存在，目前没有一套通用的泛解析记录判别方案能保证准确率，
// 建议在子域名搜集结束后，针对目标的情况定制方案对可疑子域名进一步筛选。
func (e *Engine) checker(wg *sync.WaitGroup) {
	defer func() {
		close(e.toRecorder)
		//log.Println("Checker done.")
		wg.Done()
	}()

	//wgWorker := sync.WaitGroup{}
	for t := range e.toChecker {
		if len(t.Answer) != len(e.check.dnsMark) {
			e.toRecorder <- t
			continue
		}
		//wgWorker.Add(1) //todo 限量
		//go e.doCheck(t, &wgWorker)
		t.Answer[0].Header().Name = "" // 将DNS查询的域名置空，方便比较
		for i, v := range e.check.dnsMark {
			if !dns.IsDuplicate(v, t.Answer[i]) {
				e.toRecorder <- t
				continue
			}
		}
	}
	//wgWorker.Wait()
}

//// doCheck 判断一个子域名是否有效
//func (e *Engine) doCheck(t *module.Task, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	t.Answer[0].Header().Name = "" // 将DNS查询的域名置空，方便比较
//	for i, v := range e.check.dnsMark {
//		if !dns.IsDuplicate(v, t.Answer[i]) {
//			e.toRecorder <- t
//			return
//		}
//	}
//	if e.check.existHTTP {
//		reqs, err := getHTTP("http://" + strings.TrimSuffix(t.Subdomain, "."))
//		if err != nil || len(reqs) != len(e.check.httpMark) {
//			e.toRecorder <- t
//			return
//		}
//		for i, v := range e.check.httpMark {
//			if reqs[i].status != v.status || reqs[i].url != v.url {
//				e.toRecorder <- t
//				return
//			}
//		}
//	}
//}

//
//type req struct {
//	url    string
//	status int
//}
//
//// getHTTP 记录GET访问url时的所有跳转和状态码
//func getHTTP(url string) ([]*req, error) {
//	r := &http.Client{Timeout: NETTIMEOUT} //todo 超时没起作用？
//	resp, err := r.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	resp.Body.Close()
//
//	reqs := make([]*req, 0, 10)
//	saveReq(&reqs, resp)
//	return reqs, nil
//}
//
//func saveReq(reqs *[]*req, resp *http.Response) {
//	if resp.Request.Response != nil {
//		saveReq(reqs, resp.Request.Response)
//	}
//	*reqs = append(*reqs, &req{
//		url:    resp.Request.URL.String(),
//		status: resp.StatusCode,
//	})
//}
