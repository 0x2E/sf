package wildcard

import "github.com/0x2E/sf/internal/option"

// WildcardModel 泛解析模块主体结构
type WildcardModel struct {
	mode      int
	blacklist map[string]string
	check     checkAction // 检测函数
}

// New 返回泛解析处理模块，没有初始化黑名单
func New(option *option.Option) *WildcardModel {
	var c checkAction
	switch option.Wildcard.Mode {
	case 1:
		c = checkMode1
	case 2:
		c = checkMode2
	}

	return &WildcardModel{
		mode:      option.Wildcard.Mode,
		blacklist: make(map[string]string),
		check:     c,
	}
}

// Check 返回的结果表示是否应该丢弃，下列情况下应丢弃：
// 1. 宽松模式下命中黑名单
// 2. 严格模式下命中黑名单，但无法获取网页标题
// 3. 严格模式下命中黑名单，网页标题与黑名单中的高度相似
func (w *WildcardModel) Check(subdomain, ip string) bool {
	if _, ok := w.blacklist[ip]; !ok {
		return false
	}
	res := w.check(w, subdomain, ip)
	return res
}

func (w *WildcardModel) BlacklistLen() int { return len(w.blacklist) }
