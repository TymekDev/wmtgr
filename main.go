package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

// TODO: request webmetions from URL at the start
// TODO: check every N minutes
// TODO: send telegram notification

func main() {
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
				Name:     "token-webmention, wm",
				Usage:    "webmention.io API token",
				Required: true,
			},
		},
		Commands: []cli.Command{
			{
				Name:  "fetch",
				Usage: "fetch webmentions once, print, and exit",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
			},
			{
				Name:  "relay",
				Usage: "fetch webmentions periodically and relay them to Telegram",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "token-telegram, tg",
						Usage:    "Telegram bot API token",
						Required: true,
					},
					cli.DurationFlag{
						Name:  "interval, n",
						Usage: "interval between checks",
						Value: 60 * time.Minute,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
