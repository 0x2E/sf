package controller

import (
	"fmt"
	"github.com/0x2E/sf/model"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"regexp"
)

type setApp func(app *model.App, ctx *cli.Context)

var setAppList = []setApp{
	setDomain,
	setDict,
	setThread,
	setOutput,
	setResolver,
	setQueue,
	setWildcard,
	setRetry,
}

// setDomain 设置域名
func setDomain(app *model.App, c *cli.Context) {
	re, _ := regexp.Compile(`^(?:\w*://)?((?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6})/?`)
	domain := re.FindStringSubmatch(c.String("url"))
	if len(domain) == 1 || len(domain[1]) < 4 {
		log.Fatal("invalid domain name")
	}
	app.Domain = domain[1]
}

// setDict 设置字典
func setDict(app *model.App, c *cli.Context) {
	app.Dict = c.String("dict")
	if app.Dict != "" {
		f, err := os.Open(app.Dict)
		if err != nil {
			log.Fatalf("failed to open %s: %s", app.Dict, err.Error())
		}
		f.Close()
	}
}

// setOutput 设置输出文件
func setOutput(app *model.App, c *cli.Context) {
	app.Output = c.String("output")
	if app.Output == "" {
		app.Output = fmt.Sprintf("%s-%d.txt", app.Domain, app.Start.Unix())
	}
}

// setResolver 设置DNS服务器
func setResolver(app *model.App, c *cli.Context) {
	app.Resolver = c.String("resolver") + ":53"
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	var err error
	for i := 0; i < 2; i++ { // 多试几次
		_, err = dns.Exchange(m, app.Resolver)
		if err == nil {
			break
		}
	}
	if err != nil { // 重试之后仍有错误
		log.Fatal("resolver may be invalid: ", err.Error())
	}
}

// setThread 设置并发数
func setThread(app *model.App, c *cli.Context) {
	app.Thread = c.Int("thread")
	mustPositive("thread", app.Thread)
}

// setQueue 设置fuzz任务队列长度
func setQueue(app *model.App, c *cli.Context) {
	app.Queue = c.Int("queue")
	mustPositive("queue", app.Queue)
}

// setWildcard 设置泛解析模式
func setWildcard(app *model.App, c *cli.Context) {
	app.Wildcard.Mode = c.Int("wildcardMode")
	mustPositive("wildcardMode", app.Wildcard.Mode)

	app.Wildcard.BlacklistMaxLen = c.Int("wildcardBlacklistMaxLen")
	mustPositive("wildcardBlacklistMaxLen", app.Wildcard.BlacklistMaxLen)
}

// setRetry 设置重试次数
func setRetry(app *model.App, c *cli.Context) {
	app.Retry = c.Int("retry")
	mustPositive("retry", app.Retry)
}

func mustPositive(name string, n int) {
	if n < 0 {
		log.Fatalf("`%s` cannot be negative", name)
	}
}
