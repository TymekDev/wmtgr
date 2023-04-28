package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Bot struct {
	token  string
	chatID string
}

func NewBot(token, chatID string) (*Bot, error) {
	tg := &Bot{
		token:  token,
		chatID: chatID,
	}

	name, err := tg.Name()
	if err != nil {
		return nil, err
	}
	log.Printf("INFO authorized as @%s in Telegram bot API\n", name)

	return tg, nil
}

func (tg *Bot) Name() (string, error) {
	b, err := tg.do(nil, "getMe")
	if err != nil {
		return "", err
	}

	result := struct {
		Result struct {
			Username string `json:"username"`
		} `json:"result"`
	}{}
	if err := json.Unmarshal(b, &result); err != nil {
		return "", err
	}

	return result.Result.Username, nil
}

func (tg *Bot) Send(message string) error {
	payload := struct {
		ChatID  string `json:"chat_id"`
		Message string `json:"text"`
	}{
		ChatID:  tg.chatID,
		Message: message,
	}

	if _, err := tg.do(payload, "sendMessage"); err != nil {
		return err
	}

	return nil
}

func (tg *Bot) do(payload any, method string) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, tg.url(method), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("ERROR", err)
		}

		return nil, fmt.Errorf("Telegram bot API request failed: %s: %s", resp.Status, string(b))
	}

	return io.ReadAll(resp.Body)
}

func (b *Bot) url(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.token, method)
}
