package engine

import (
	"github.com/0x2E/sf/internal/module"
	"github.com/miekg/dns"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/ratelimit"
	"log"
	"net"
	"sync"
	"time"
)

func (e *Engine) enumerator(wg *sync.WaitGroup) {
	defer func() {
		//log.Println("Enumerator done.")
		wg.Done()
	}()

	taskMap := &sync.Map{}
	toSender := make(chan *module.Task, e.conf.Thread)
	wgWorker := sync.WaitGroup{}
	var toNext chan *module.Task
	if e.check.existWildcard {
		toNext = e.toChecker
	} else {
		toNext = e.toRecorder
	}
	defer close(toNext)

	for i := 0; i < e.conf.Thread; i++ {
		udpConn, err := net.Dial("udp", e.conf.Resolver)
		if err != nil {
			//log.Println(err)
			continue
		}
		worker := &enumerateWorker{
			conn:       &dns.Conn{Conn: udpConn},
			bar:        e.bar,
			senderDone: false,
		}
		go worker.sender(toSender)
		wgWorker.Add(1)
		go worker.receiver(&wgWorker, toNext, taskMap)
	}

	go func() {
		rl := ratelimit.New(e.conf.Rate)
		var ok bool
		for t := range e.toEnumerator {
			if _, ok = taskMap.Load(t.Subdomain); ok {
				continue
			}
			rl.Take()
			toSender <- t
			taskMap.Store(t.Subdomain, t)
		}
		// 将超时的子域名再次加入任务队列
		for i := 0; i < e.conf.Retry; i++ {
			count := 0
			taskMap.Range(func(s, t interface{}) bool {
				if t, ok := t.(*module.Task); ok {
					if t.Received || time.Now().Unix()-t.Time < int64(NETTIMEOUT.Seconds()) {
						return true
					}
					rl.Take()
					toSender <- t
					count++
				}
				return true
			})
			if count == 0 {
				break
			}
		}
		close(toSender)
	}()

	wgWorker.Wait()
}

type enumerateWorker struct {
	conn       *dns.Conn
	bar        *progressbar.ProgressBar
	senderDone bool
}

func (ew *enumerateWorker) sender(toSender <-chan *module.Task) {
	var msg *dns.Msg
	var err error
	for t := range toSender {
		msg = (&dns.Msg{}).SetQuestion(t.Subdomain, dns.TypeA)
		ew.conn.SetWriteDeadline(time.Now().Add(NETTIMEOUT))
		err = ew.conn.WriteMsg(msg)
		if err != nil {
			log.Println(err)
			continue
		}
		t.Time = time.Now().Unix()
	}
	ew.senderDone = true
}

func (ew *enumerateWorker) receiver(wg *sync.WaitGroup, toNext chan<- *module.Task, taskMap *sync.Map) {
	defer func() {
		wg.Done()
		ew.conn.Close()
	}()

	var ttmp interface{}
	var t *module.Task
	var ok bool
	for {
		ew.conn.SetReadDeadline(time.Now().Add(NETTIMEOUT))
		msg, err := ew.conn.ReadMsg()
		if err != nil {
			if err2, ok := err.(net.Error); ok && err2.Timeout() && ew.senderDone {
				// 在senderDone的情况下读超时，可以认为没有待接收的相应包，直接退出
				// todo 考虑更多异常情况
				return
			}
			continue
		}
		ew.bar.Add(1)
		if ttmp, ok = taskMap.Load(msg.Question[0].Name); ok {
			if t, ok = ttmp.(*module.Task); ok {
				t.Received = true
				if msg.Rcode != dns.RcodeSuccess || len(msg.Answer) == 0 {
					continue
				}
				t.Answer = msg.Answer
				toNext <- t
			}
		}
	}
}
