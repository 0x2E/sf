package fuzz

import (
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/internal/util/dnsudp"
	"net"
	"strings"
	"sync"
)

type udpContext struct {
	queue        chan struct{} // 任务队列，让发送和接收的速度有所同步，防止发送过快但接收处理跟不上而冲爆缓存，也是防止接收者中的数据处理goroutine无限增多
	receiverDone chan struct{}
	recorderDone chan struct{}
	recordQueue  chan string
	receivedMap  map[string]struct{}
	resultMap    map[string]string
	mu           sync.Mutex
}

// consumer
func consumer(f *FuzzModule, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("udp", f.option.resolver)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx := &udpContext{
		queue:        make(chan struct{}, f.option.queue),
		receiverDone: make(chan struct{}),
		recorderDone: make(chan struct{}),
		recordQueue:  make(chan string, f.option.queue),
		receivedMap:  make(map[string]struct{}),
		resultMap:    make(map[string]string),
	}

	go receiver(conn, ctx, f)
	go recorder(ctx)

	// 发送UDP
	var id uint16 = 1 // TODO 发送量可能超过uint16的上限
	for entry := range ch {
		id++
		err := dnsudp.Send(conn, entry, id, rawdns.QTypeA)
		if err != nil {
			continue
		}
		ctx.queue <- struct{}{}
	}
	close(ctx.queue)

	<-ctx.receiverDone
	<-ctx.recorderDone

	// 将本消费者接收到的子域名从模块结构中删去
	tmp := make([]string, 0, 5000)
	f.unReceived.mu.Lock()
	for k := range f.unReceived.data {
		if _, ok := ctx.receivedMap[f.unReceived.data[k]]; ok {
			continue
		}
		tmp = append(tmp, f.unReceived.data[k])
	}
	f.unReceived.data = tmp
	f.unReceived.mu.Unlock()

	// 将本消费者的结果汇总到模块结构中
	f.result.mu.Lock()
	for k := range ctx.resultMap {
		f.result.data[k] = ctx.resultMap[k]
	}
	f.result.mu.Unlock()
}

// receiver UDP接收者
func receiver(conn net.Conn, ctx *udpContext, f *FuzzModule) {
	wg := sync.WaitGroup{} // 用于处理数据的goroutine
	for range ctx.queue {
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
	close(ctx.recordQueue)
	ctx.receiverDone <- struct{}{}
}

// handleResp 处理接收到的DNS响应
func handleResp(resp *rawdns.Message, wg *sync.WaitGroup, ctx *udpContext, f *FuzzModule) {
	defer wg.Done()

	if len(resp.Questions) == 0 {
		return
	}
	ctx.recordQueue <- resp.Questions[0].QNAME // 标记已收到此子域名的响应

	if len(resp.Answers) == 0 {
		return
	}
	subdomain := resp.Questions[0].QNAME
	ip := net.IP(resp.Answers[len(resp.Answers)-1].RDATA).String() // 用切片的最后一个，前面的可能是CNAME

	if ok := f.wildcard.Check(subdomain, ip); !ok {
		ctx.mu.Lock()
		ctx.resultMap[subdomain] = ip
		ctx.mu.Unlock()
	}
}

// recorder 记录已接收到的子域名查询结果
func recorder(ctx *udpContext) {
	for domain := range ctx.recordQueue {
		ctx.receivedMap[domain] = struct{}{}
	}
	ctx.recorderDone <- struct{}{}
}
