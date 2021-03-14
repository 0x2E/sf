package controller

import (
	"errors"
	"fmt"
	"github.com/0x2E/rawdns"
	"github.com/0x2E/sf/model"
	"github.com/0x2E/sf/util/dnsudp"
	"github.com/urfave/cli/v2"
	"net"
	"os"
	"regexp"
)

type setApp func(app *model.App, ctx *cli.Context) error

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
func setDomain(app *model.App, c *cli.Context) error {
	input := c.String("url")
	re, _ := regexp.Compile(`^(?:\w*://)?((?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6})/?`)
	domain := re.FindStringSubmatch(input)
	if len(domain) == 1 {
		return errors.New("cannot resolve to a domain name")
	}
	if len(domain[1]) < 4 {
		return errors.New("bad domain name format" + domain[1])
	}

	app.Domain = domain[1]
	return nil
}

// setDict 设置字典
func setDict(app *model.App, c *cli.Context) error {
	input := c.String("dict")
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	f.Close()

	app.Dict = input
	return nil
}

// setOutput 设置输出文件
func setOutput(app *model.App, c *cli.Context) error {
	var outFile string
	input := c.String("output")
	if input != "" {
		outFile = input
		return nil
	} else {
		outFile = fmt.Sprintf("%s-%d.txt", app.Domain, app.Start.Unix())
	}

	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	f.Close()
	app.Output = outFile
	return nil
}

// setResolver 设置DNS服务器
func setResolver(app *model.App, c *cli.Context) error {
	input := c.String("resolver")
	app.Resolver = input + ":53"

	// 测试是否可用
	conn, err := net.Dial("udp", app.Resolver)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot create socket to test resolver [%s]: %s\n", app.Resolver, err.Error()))
	}
	defer conn.Close()

	if err := dnsudp.Send(conn, "google.com", 123, rawdns.QTypeA); err != nil {
		return errors.New(fmt.Sprintf("cannot send DNS udp to resolver [%s]: %s\n", app.Resolver, err.Error()))
	}

	if _, err := dnsudp.Receive(conn, 2); err != nil {
		return errors.New(fmt.Sprintf("cannot receive DNS udp from resolver [%s]: %s\n", app.Resolver, err.Error()))
	}

	return nil
}

// setThread 设置并发数
func setThread(app *model.App, c *cli.Context) error {
	input := c.Int("thread")
	if input < 0 || input > 999999 {
		return errors.New("the number of thread must be between 0 and 999999")
	}
	app.Thread = input
	return nil
}

// setQueue 设置UDP请求-发送队列长度
func setQueue(app *model.App, c *cli.Context) error {
	input := c.Int("queue")
	if input < 0 || input > 99999 {
		return errors.New("the number of retries must be between 0 and 99999")
	}
	app.Queue = input
	return nil
}

// setWildcard 设置泛解析模式
func setWildcard(app *model.App, c *cli.Context) error {
	input := c.Int("wildcard")
	if input != 1 && input != 2 {
		return errors.New("the parameter wildcard(-d) must be 1 or 2")
	}
	app.Wildcard = input
	return nil
}
