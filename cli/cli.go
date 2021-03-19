package cli

import (
	"github.com/0x2E/sf/controller"
	"github.com/urfave/cli/v2"
	"log"
)

const (
	threadDefault   = 300
	queueDefault    = 300
	wildcardDefault = 1
	dictDefault     = "./dict.txt"
	resolverDefault = "8.8.8.8"
)

// Handle 处理cli输入
func Handle(args []string) {
	c := &cli.App{
		Name:    "sf",
		Usage:   "subdomain finder - https://github.com/0x2E/sf",
		Version: "v0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "url", Aliases: []string{"u"}, Usage: "target url or domain name", Required: true},
			&cli.StringFlag{Name: "dict", Aliases: []string{"d"}, Usage: "load dictionary from `FILE`", Value: dictDefault},
			&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "output results to `FILE`"},
			&cli.StringFlag{Name: "resolver", Aliases: []string{"r"}, Usage: "DNS resolver", Value: resolverDefault},
			&cli.IntFlag{Name: "thread", Aliases: []string{"t"}, Usage: "the number of concurrent, each will occupy a temporary port of the system", Value: threadDefault},
			&cli.IntFlag{Name: "queue", Aliases: []string{"q"}, Usage: "the size of the udp sending and receiving queues. it depends on your system network conditions. the higher the faster, but it is easy to cause omissions.", Value: queueDefault},
			&cli.IntFlag{Name: "wildcard", Aliases: []string{"w"}, Usage: "modes for handling wildcard DNS", Value: wildcardDefault},
		},
		Action: controller.Handle,
		//UseShortOptionHandling: true,
		HideVersion:     true,
		HideHelpCommand: true,
	}
	err := c.Run(args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
