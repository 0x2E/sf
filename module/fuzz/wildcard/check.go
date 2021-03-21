package wildcard

import (
	"github.com/antlabs/strsim"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// checkAction 泛解析黑名单检测函数
type checkAction func(w *WildcardModel, subdomain, ip string) bool

// checkMod1 宽松模式测试
// 之前已经判断是否匹配到黑名单了，这里没有继续的操作，留空以备扩展
func checkMod1(w *WildcardModel, subdomain, ip string) bool { return true }

// checkMod2 严格模式测试
func checkMod2(w *WildcardModel, subdomain, ip string) bool {
	title, err := getPageTitle("http" + subdomain)
	if err != nil {
		// 无法获取标题，丢弃
		return true
	}
	rank := strsim.Compare(title, w.Blacklist[ip])
	if rank > 0.5 {
		// 相似度较高，丢弃
		return true
	}
	return false
}

// getPageTitle 获取url的网页标题，2秒超时。无法获取到、获取到的都是空格换行符时，返回空字符串
func getPageTitle(url string) (string, error) {
	c := http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	re, _ := regexp.Compile(`<title>([\s\S]*)</title>`)
	title := re.FindStringSubmatch(string(body))
	if len(title) < 2 {
		return "", err
	}
	return strings.TrimSpace(title[1]), nil
}

// randString 生成长度为n的随机字符串
func randString(n int) string {
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
