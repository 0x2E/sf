package wildcard

import (
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/internal/util/dnsudp"
	"net"
	"strings"
	"sync"
)

// InitBlacklist 初始化黑名单。不放在New中一起初始化，是为了避免初始化黑名单时间过长导致命令行卡住没有输出
func (w *WildcardModel) InitBlacklist(domain, resolver string, queueLen, mode, maxlen int) {
	conn, err := net.Dial("udp", resolver)
	if err != nil {
		w.blacklist = make(map[string]string)
		return
	}
	defer conn.Close()

	blacklist := &sync.Map{}
	queue := make(chan struct{}, queueLen)
	receiverDone := make(chan struct{})

	go receiver(conn, queue, mode, blacklist, receiverDone)

	for i := 0; i < maxlen; i++ {
		domain := randString(12) + "." + domain
		if err := dnsudp.Send(conn, domain, uint16(i), rawdns.QTypeA); err != nil {
			continue
		}
		queue <- struct{}{}
	}
	close(queue)

	<-receiverDone

	blacklist.Range(func(k, v interface{}) bool {
		w.blacklist[k.(string)] = v.(string)
		return true
	})
}

//receiver 接收爆破黑名单时返回的报文并保存
func receiver(conn net.Conn, queue <-chan struct{}, mode int, blacklist *sync.Map, done chan struct{}) {
	wg := sync.WaitGroup{}

	var recordOne func(wg *sync.WaitGroup, subdomain, ip string, blacklist *sync.Map)
	switch mode {
	case 1:
		// 宽松模式：仅记录IP
		recordOne = func(wg *sync.WaitGroup, subdomain, ip string, blacklist *sync.Map) {
			defer wg.Done()
			blacklist.Store(ip, "")
		}
	case 2:
		// 严格模式：记录IP和网页标题
		recordOne = func(wg *sync.WaitGroup, subdomain, ip string, blacklist *sync.Map) {
			defer wg.Done()

			title, _ := getPageTitle("http://" + subdomain)
			// 即使是空字符串也要保存
			blacklist.Store(ip, title)
		}
	}

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
		if _, ok := blacklist.Load(ip); ok {
			// 该IP已存在于黑名单，跳过
			continue
		}

		wg.Add(1)
		go recordOne(&wg, resp.Questions[0].QNAME, ip, blacklist) // 因为可预测黑名单长度不会太大，所以没必要限制goroutine的数量
	}
	wg.Wait()
	done <- struct{}{}
}
