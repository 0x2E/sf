package wildcard

import (
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/util/dnsudp"
	"github.com/pkg/errors"
	"net"
	"sync"
)

// WildcardModel 泛解析模块主体结构
type WildcardModel struct {
	mode      int
	blacklist struct {
		maxLen int
		data   map[string]string
	}
	c checkAction // 检测函数
	b blAction    // 黑名单初始化函数
}

// New 返回一个泛解析结构体
func New() *WildcardModel {
	return &WildcardModel{
		blacklist: struct {
			maxLen int
			data   map[string]string
		}{data: make(map[string]string)},
	}
}

// Init 初始化黑名单，设置检测函数
func (w *WildcardModel) Init(app *model.App) error {
	w.mode = app.Wildcard.Mode
	w.blacklist.maxLen = app.Wildcard.BlacklistMaxLen
	switch w.mode {
	case 1:
		w.c = checkMode1
		w.b = blMode1
	case 2:
		w.c = checkMode2
		w.b = blMode2
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

	go blReceiver(conn, queue, blacklist, w, receiverDone)

	for i := 0; i < w.blacklist.maxLen; i++ {
		domain := randString(12) + "." + app.Domain
		if err := dnsudp.Send(conn, domain, uint16(i), rawdns.QTypeA); err != nil {
			continue
		}
		queue <- struct{}{}
	}
	close(queue)

	<-receiverDone

	blacklist.Range(func(k, v interface{}) bool {
		w.blacklist.data[k.(string)] = v.(string)
		return true
	})
	return nil
}

// Check 返回的结果表示是否应该丢弃，下列情况下应丢弃：
// 1. 宽松模式下命中黑名单
// 2. 严格模式下命中黑名单，但无法获取网页标题
// 3. 严格模式下命中黑名单，网页标题与黑名单中的高度相似
func (w *WildcardModel) Check(subdomain, ip string) bool {
	if _, ok := w.blacklist.data[ip]; !ok {
		return false
	}
	res := w.c(w, subdomain, ip)
	return res
}

func (w *WildcardModel) BlacklistLen() int { return len(w.blacklist.data) }
