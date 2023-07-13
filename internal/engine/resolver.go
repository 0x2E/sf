package engine

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/module"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
)

func (e *Engine) resolver(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.WithField("step", "resolver")

	var toNext chan *module.Task
	if e.needCheck {
		toNext = e.toChecker
	} else {
		toNext = e.toRecorder
	}
	defer close(toNext)

	taskMap := &sync.Map{}
	toSender := make(chan *module.Task, conf.C.Concurrent)
	wgWorker := sync.WaitGroup{}
	sendCount := &atomic.Uint64{}
	recvCount := &atomic.Uint64{}

	for i := 0; i < conf.C.Concurrent; i++ {
		// TODO: dns.DialWithTLS
		udpConn, err := net.Dial("udp", conf.C.Resolver)
		if err != nil {
			logger.Error(err)
			continue
		}
		wgWorker.Add(1)
		worker := &resolverWorker{
			conn:       &dns.Conn{Conn: udpConn},
			sendCount:  sendCount,
			recvCount:  recvCount,
			logger:     logger.WithField("concurrent_id", i),
			senderDone: false,
		}
		go worker.sender(toSender)
		go worker.receiver(&wgWorker, toNext, taskMap)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(conf.C.StatisticsInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				logger.Infof("[statistics] send %d, recv %d", sendCount.Load(), recvCount.Load())
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		rl := ratelimit.New(conf.C.Rate)
		for t := range e.toResolver {
			if _, ok := taskMap.Load(t.DomainName); ok {
				continue
			}
			rl.Take()
			toSender <- t
			taskMap.Store(t.DomainName, t)
		}

		// retry failed items
		for i := 0; i < conf.C.Retry; i++ {
			retryList := make([]*module.Task, 0)
			taskMap.Range(func(s, t interface{}) bool {
				if t, ok := t.(*module.Task); ok {
					if t.Received || time.Now().Unix()-t.LastQueryAt < int64(NetTimeout.Seconds()) {
						return true
					}
					retryList = append(retryList, t)
				}
				return true
			})
			if len(retryList) == 0 {
				break
			}

			logger.Infof("#%d retry: %d items", i+1, len(retryList))
			for _, t := range retryList {
				rl.Take()
				toSender <- t
			}
		}
		close(toSender)
	}()

	wgWorker.Wait()
}

type resolverWorker struct {
	conn       *dns.Conn
	sendCount  *atomic.Uint64
	recvCount  *atomic.Uint64
	senderDone bool

	logger *logrus.Entry
}

func (w *resolverWorker) sender(toSender <-chan *module.Task) {
	defer func() {
		w.senderDone = true
	}()

	for t := range toSender {
		msg := (&dns.Msg{}).SetQuestion(t.DomainName, dns.TypeA)
		w.conn.SetWriteDeadline(time.Now().Add(NetTimeout))
		if err := w.conn.WriteMsg(msg); err != nil {
			w.logger.Debug(err)
			continue
		}
		t.LastQueryAt = time.Now().Unix()
		w.sendCount.Add(1)
	}
}

func (w *resolverWorker) receiver(wg *sync.WaitGroup, toNext chan<- *module.Task, taskMap *sync.Map) {
	defer func() {
		wg.Done()
		w.conn.Close()
	}()

	for {
		w.conn.SetReadDeadline(time.Now().Add(NetTimeout))
		msg, err := w.conn.ReadMsg()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				if w.senderDone {
					// if receiver timeout after sender done, we assumed there are no more incoming packets
					return
				}
				continue
			}
			w.logger.Debug(err)
			continue
		}
		if t, ok := taskMap.Load(msg.Question[0].Name); ok {
			if task, ok := t.(*module.Task); ok {
				w.recvCount.Add(1)
				task.Received = true
				if msg.Rcode != dns.RcodeSuccess || len(msg.Answer) == 0 {
					continue
				}
				task.Answer = msg.Answer
				toNext <- task
			}
		}
	}
}
