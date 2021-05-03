package controller

import (
	"github.com/urfave/cli/v2"
	"log"
)

// Cli 处理cli输入
func Cli(args []string) {
	c := &cli.App{
		Name:    "sf",
		Usage:   "subdomain finder - https://github.com/0x2E/sf",
		Version: "v0.1.1",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "url", Aliases: []string{"u"}, Usage: "Target url or domain name", Required: true},
			&cli.StringFlag{Name: "dict", Aliases: []string{"d"}, Usage: "Load dictionary from a file"},
			&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Output results to a file"},
			&cli.StringFlag{Name: "resolver", Aliases: []string{"r"}, Usage: "[fuzz] DNS resolver", Value: "8.8.8.8"},
			&cli.IntFlag{Name: "thread", Aliases: []string{"t"}, Usage: "[fuzz] The number of threads. Each thread will occupy a temporary port of the system until the end of the fuzz", Value: 100},
			&cli.IntFlag{Name: "queue", Aliases: []string{"q"}, Usage: "[fuzz] The length of the task queue. Too high may fill the system socket buffer and cause packet loss", Value: 100},
			&cli.IntFlag{Name: "wildcardMode", Aliases: []string{"w"}, Usage: "[fuzz] Two modes (1 or 2) for processing wildcard records. Mode 1 is only based on the IP blacklist. Mode 2 matches the IP blacklist, compares the similarity of web page titles after hits, and degenerates to mode 1 if port 80 cannot be accessed.", Value: 1},
			&cli.IntFlag{Name: "wildcardBlacklistMaxLen", Aliases: []string{"wl"}, Usage: "[fuzz] The maximum length of the IP blacklist for wildcard records", Value: 1000},
			&cli.IntFlag{Name: "retry", Aliases: []string{"R"}, Usage: "[fuzz] The number of retries", Value: 2},
		},
		Action: handle,
		//UseShortOptionHandling: true,
		HideVersion:     true,
		HideHelpCommand: true,
	}
	err := c.Run(args)
	if err != nil {
		log.Println(err.Error())
	}
}
