package controller

import (
	"fmt"
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/util/dnsudp"
	"github.com/urfave/cli/v2"
	"log"
	"net"
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
	var outFile string
	input := c.String("output")
	if input != "" {
		outFile = input
		return
	} else {
		outFile = fmt.Sprintf("%s-%d.txt", app.Domain, app.Start.Unix())
	}

	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("failed to set output file: ", err.Error())
	}
	f.Close()

	app.Output = outFile
}

// setResolver 设置DNS服务器
func setResolver(app *model.App, c *cli.Context) {
	input := c.String("resolver")
	app.Resolver = input + ":53"

	// 测试是否可用
	conn, err := net.Dial("udp", app.Resolver)
	if err != nil {
		log.Fatalf("cannot create socket to test resolver [%s]: %s\n", app.Resolver, err.Error())
	}
	defer conn.Close()

	if err := dnsudp.Send(conn, "google.com", 123, rawdns.QTypeA); err != nil {
		log.Fatalf("cannot send DNS udp to resolver [%s]: %s\n", app.Resolver, err.Error())
	}

	if _, err := dnsudp.Receive(conn, 2); err != nil {
		log.Fatalf("cannot receive DNS udp from resolver [%s]: %s\n", app.Resolver, err.Error())
	}
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
