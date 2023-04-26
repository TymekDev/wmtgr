package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

// TODO: request webmetions from URL at the start
// TODO: check every N minutes
// TODO: send telegram notification

func main() {
	var (
		tokenWebmention string
		tokenTelegram   string
		interval        time.Duration
	)

	app := &cli.App{
		Name:        "wmtgr",
		Usage:       "webmentions to telegram relay",
		Description: "wmtgr periodically checks webmention.io for new webmentions and sends them to Telegram using Telegram bot API.",
		Authors: []cli.Author{
			{
				Name:  "Tymoteusz Makowski",
				Email: "tymek.makowski@gmail.com",
			},
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "token-webmention, wm",
				Usage:       "webmention.io API token",
				Required:    true,
				Destination: &tokenWebmention,
			},
			cli.StringFlag{
				Name:        "token-telegram, tg",
				Usage:       "Telegram bot API token",
				Destination: &tokenTelegram,
			},
			cli.DurationFlag{
				Name:        "interval, n",
				Usage:       "interval between checks",
				Value:       60 * time.Minute,
				Destination: &interval,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
