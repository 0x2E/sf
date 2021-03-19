package wildcard

import (
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/util/dnsudp"
	"github.com/pkg/errors"
	"net"
	"sync"
)

type WildcardModel struct {
	Mod       int
	Blacklist map[string]string // 黑名单
	c         checkAction       // 检测函数
	b         blAction          // 黑名单初始化函数
}

// Init 初始化黑名单，设置检测函数
func (w *WildcardModel) Init(app *model.App) error {
	// 选择检测函数
	w.Mod = app.Wildcard
	switch w.Mod {
	case 1:
		w.c = checkMod1
		w.b = blMod1
	case 2:
		w.c = checkMod2
		w.b = blMod2
	}

	// 初始化黑名单
	conn, err := net.Dial("udp", app.Resolver)
	if err != nil {
		return errors.Wrap(err, "failed to create socket")
	}
	defer conn.Close()

	blacklist := &sync.Map{}
	queue := make(chan struct{}, app.Queue)
	receiverDone := make(chan struct{})

	go blThread(conn, queue, blacklist, w, receiverDone)

	for i := 0; i < 6000; i++ { // 经测试，6000次基本可以爆破出所有泛解析目的IP
		domain := randString(12) + "." + app.Domain
		if err := dnsudp.Send(conn, domain, uint16(i), rawdns.QTypeA); err != nil {
			continue
		}
		queue <- struct{}{}
	}
	close(queue)

	<-receiverDone

	blacklist.Range(func(k, v interface{}) bool {
		w.Blacklist[k.(string)] = v.(string)
		return true
	})
	return nil
}

// Check 返回的结果表示是否应该丢弃，下列情况下应丢弃：
// 1. 宽松模式下命中黑名单
// 2. 严格模式下命中黑名单，但无法获取网页标题
// 3. 严格模式下命中黑名单，网页标题与黑名单中的高度相似
func (w *WildcardModel) Check(subdomain, ip string) bool {
	if _, ok := w.Blacklist[ip]; !ok {
		return false
	}
	res := w.c(w, subdomain, ip)
	return res
}
