package fuzz

import (
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/util/dnsudp"
	"net"
	"strings"
	"sync"
)

type udpContext struct {
	queue        chan struct{} // 任务队列，让发送和接收的速度有所同步，防止发送过快但接收处理跟不上而冲爆缓存，也是防止接收者中的数据处理goroutine无限增多
	receiverDone chan struct{}
	result       map[string]string
	mu           sync.Mutex
}

// consumer
func consumer(ch <-chan string, wg *sync.WaitGroup, app *model.App, f *FuzzModule) {
	defer wg.Done()

	conn, err := net.Dial("udp", app.Resolver)
	if err != nil {
		return
	}
	defer conn.Close()

	//TODO retry

	ctx := &udpContext{
		queue:        make(chan struct{}, app.Queue),
		receiverDone: make(chan struct{}),
		result:       make(map[string]string),
	}

	go receiver(conn, ctx, f) // udp接收者

	// 发送UDP
	var id uint16 = 1 // TODO 发送量可能超过uint16的上限
	for entry := range ch {
		id++
		subdomain := entry + "." + app.Domain
		err := dnsudp.Send(conn, subdomain, id, rawdns.QTypeA)
		if err != nil {
			continue
		}
		ctx.queue <- struct{}{}
	}
	close(ctx.queue)

	<-ctx.receiverDone

	f.Result.Mu.Lock()
	for k := range ctx.result {
		f.Result.Data[k] = ctx.result[k]
	}
	f.Result.Mu.Unlock()
}

// receiver UDP接收者
func receiver(conn net.Conn, ctx *udpContext, f *FuzzModule) {
	wg := sync.WaitGroup{} // 用于处理数据的goroutine
	for range ctx.queue {
		//TODO 有接收数小于发送数的情况，1.缓存满后被丢弃？ 2.某些请求被链路或DNS服务器丢弃？
		resp, err := dnsudp.Receive(conn, 2)
		if err != nil {
			// 只要等待的时间足够长，依旧收不到任何响应包，就可以认为此时queue中代表的请求都无法接收到回复了，所以清空queue
			// 超时错误内容：read udp 192.168.0.102:54012->8.8.8.8:53: i/o timeout
			if strings.Contains(err.Error(), "timeout") {
				l := len(ctx.queue)
				for i := 0; i < l; i++ {
					<-ctx.queue
				}
			}
			continue
		}

		// 处理接收的数据 TODO 对最大goroutine数做限制
		wg.Add(1)
		go handleResp(resp, &wg, ctx, f)
	}

	wg.Wait()
	ctx.receiverDone <- struct{}{}
}

// handleResp 处理接收到的DNS响应
func handleResp(resp *rawdns.Message, wg *sync.WaitGroup, ctx *udpContext, f *FuzzModule) {
	defer wg.Done()

	if len(resp.Answers) == 0 || len(resp.Questions) == 0 {
		return
	}

	subdomain := resp.Questions[0].QNAME
	ip := net.IP(resp.Answers[len(resp.Answers)-1].RDATA).String() // 用切片的最后一个，前面的可能是CNAME
	if ok := f.Wildcard.Check(subdomain, ip); !ok {
		ctx.mu.Lock()
		ctx.result[subdomain] = ip
		ctx.mu.Unlock()
	}
}
