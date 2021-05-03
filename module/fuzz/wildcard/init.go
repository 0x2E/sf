package wildcard

import (
	"github.com/0x2E/sf/util/dnsudp"
	"net"
	"strings"
	"sync"
)

type blAction func(wg *sync.WaitGroup, subdomain, ip string, blacklist *sync.Map)

// blMode1 宽松模式黑名单初始化：仅记录IP
func blMode1(wg *sync.WaitGroup, subdomain, ip string, blacklist *sync.Map) {
	defer wg.Done()
	blacklist.Store(ip, "")
}

// blMode2 严格模式黑名单初始化：记录IP和网页标题
func blMode2(wg *sync.WaitGroup, subdomain, ip string, blacklist *sync.Map) {
	defer wg.Done()

	title, _ := getPageTitle("http://" + subdomain)
	// 即使是空字符串也要保存
	blacklist.Store(ip, title)
}

//blReceiver 黑名单爆破的接收者
func blReceiver(conn net.Conn, queue <-chan struct{}, blacklist *sync.Map, w *WildcardModel, done chan struct{}) {
	wgResp := sync.WaitGroup{} // 用于判定是否该加入黑名单的goroutine
	for range queue {
		resp, err := dnsudp.Receive(conn, 2)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				l := len(queue)
				for i := 0; i < l; i++ {
					<-queue
				}
			}
			continue
		}
		if len(resp.Answers) == 0 { // 或是通过FLAGS中的RCODE、ANCOUNT等判断
			continue
		}

		ip := net.IP(resp.Answers[len(resp.Answers)-1].RDATA).String()
		if _, ok := blacklist.Load(ip); ok { // 该IP已存在于黑名单，跳过
			continue
		}

		wgResp.Add(1)
		go w.b(&wgResp, resp.Questions[0].QNAME, ip, blacklist) // 因为可预测黑名单中的IP不会太多，所以没必要限制goroutine的数量
	}
	wgResp.Wait()
	done <- struct{}{}
}
