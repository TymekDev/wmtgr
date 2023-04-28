package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
				Name:     "webmention-token, wm",
				Usage:    "webmention.io API token",
				Required: true,
			},
		},
		Commands: []cli.Command{
			{
				Name:  "fetch",
				Usage: "fetch webmentions once, print, and exit",
				Action: func(c *cli.Context) error {
					b, err := fetch(c.GlobalString("webmention-token"))
					if err != nil {
						return err
					}

					if _, err := bytes.NewReader(b).WriteTo(os.Stdout); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "relay",
				Usage: "fetch webmentions periodically and relay them to Telegram",
				Action: func(c *cli.Context) error {
					b, err := fetch(c.GlobalString("webmention-token"))
					if err != nil {
						return err
					}

					result := &struct {
						Links []any `json:"links"`
					}{}

					if err := json.Unmarshal(b, result); err != nil {
						return err
					}

					log.Println("INFO starting relay, found", len(result.Links), "webmention(s)")
					for range time.Tick(c.Duration("interval")) {
						nLast := len(result.Links)

						b, err := fetch(c.GlobalString("webmention-token"))
						if err != nil {
							log.Println("ERROR", err)
							continue
						}

						if err := json.Unmarshal(b, result); err != nil {
							log.Println("ERROR", err)
							continue
						}

						if n := nLast - len(result.Links); n > 0 {
							log.Println("INFO found", n, "new webmention(s)")
						}
					}

					return nil
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "telegram-token, tg",
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

func fetch(token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "https://webmention.io/api/mentions", strings.NewReader(fmt.Sprintf("token=%s", token)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
