package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

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
					b, err := fetch(c.GlobalString("webmention-token"), 0)
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
					_, sinceID, err := fetchAndParse(c.GlobalString("webmention-token"), 0)
					if err != nil {
						return err
					}

					log.Println("INFO starting relay")
					for range time.Tick(c.Duration("interval")) {
						sentences, id, err := fetchAndParse(c.GlobalString("webmention-token"), sinceID)
						if err != nil {
							log.Println("ERROR", err)
							continue
						}
						sinceID = id

						if n := len(sentences); n > 0 {
							log.Println("INFO found", n, "new webmention(s)")
							const sep = "\n - "
							message := fmt.Sprintf("Found %d new webmention(s):%s%s", n, sep, strings.Join(sentences, sep))
							if err := sendTelegramMessage(c.String("telegram-token"), c.String("telegram-chat-id"), message); err != nil {
								log.Println("ERROR", err)
								continue
							}
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
					cli.StringFlag{
						Name:     "telegram-chat-id, cid",
						Usage:    "Telegram chat ID",
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

type Response struct {
	Links []struct {
		ID       int `json:"id"`
		Activity struct {
			Sentence string `json:"sentence"`
		} `json:"activity"`
	} `json:"links"`
}

func fetchAndParse(token string, sinceID int) ([]string, int, error) {
	b, err := fetch(token, sinceID)
	if err != nil {
		return nil, 0, err
	}

	var resp Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, 0, err
	}

	id := sinceID
	var sentences []string
	for _, link := range resp.Links {
		if link.ID > id {
			id = link.ID
		}
		sentences = append(sentences, link.Activity.Sentence)
	}

	return sentences, id, nil
}

func fetch(token string, sinceID int) ([]byte, error) {
	uv := url.Values{
		"token":    []string{token},
		"since_id": []string{strconv.Itoa(sinceID)},
	}

	req, err := http.NewRequest(http.MethodGet, "https://webmention.io/api/mentions", strings.NewReader(uv.Encode()))
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

func sendTelegramMessage(token string, chat_id, message string) error {
	payload := struct {
		ChatID  string `json:"chat_id"`
		Message string `json:"text"`
	}{
		ChatID:  chat_id,
		Message: message,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("ERROR", err)
		}

		return fmt.Errorf("failed to send message: %s: %s", resp.Status, string(b))
	}

	return nil
}
