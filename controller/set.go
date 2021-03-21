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
}

// setDomain 设置域名
func setDomain(app *model.App, c *cli.Context) {
	input := c.String("url")
	re, _ := regexp.Compile(`^(?:\w*://)?((?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6})/?`)
	domain := re.FindStringSubmatch(input)
	if len(domain) == 1 || len(domain[1]) < 4 {
		log.Fatal("invalid domain name")
	}

	app.Domain = domain[1]
}

// setDict 设置字典
func setDict(app *model.App, c *cli.Context) {
	input := c.String("dict")
	f, err := os.Open(input)
	if err != nil {
		log.Fatal("failed to open dict file: ", err.Error())
	}
	f.Close()

	app.Dict = input
}

// setOutput 设置输出文件
func setOutput(app *model.App, c *cli.Context) {
	input := c.String("output")
	if input == "" {
		app.Output = fmt.Sprintf("%s-%d.txt", app.Domain, app.Start.Unix())
		return
	}
	app.Output = input
}

// setResolver 设置DNS服务器
func setResolver(app *model.App, c *cli.Context) {
	input := c.String("resolver") + ":53"

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	var r *dns.Msg
	var err error
	for i := 0; i < 2; i++ { // 重试几次
		r, err = dns.Exchange(m, input)
		if err == nil {
			break
		}
	}
	if err != nil || r.Rcode != dns.RcodeSuccess { // 重试之后仍有错误
		log.Fatal("resolver may be invalid: ", err.Error())
	}

	app.Resolver = input
}

// setThread 设置并发数
func setThread(app *model.App, c *cli.Context) {
	input := c.Int("thread")
	if input < 0 || input > 99999 {
		log.Fatal("thread must be between 0 and 99999")
	}

	app.Thread = input
}

// setQueue 设置UDP请求-发送队列长度
func setQueue(app *model.App, c *cli.Context) {
	input := c.Int("queue")
	if input < 0 || input > 99999 {
		log.Fatal("queue must between 0 and 99999")
	}

	app.Queue = input
}

// setWildcard 设置泛解析模式
func setWildcard(app *model.App, c *cli.Context) {
	input := c.Int("wildcard")
	if input != 1 && input != 2 {
		log.Fatal("the wildcard mod must be 1 or 2")
	}

	app.Wildcard = input
}
