package telegram

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"time"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var telegramErrMsg = `
%s

Command line flags will override environment variables

`

type Telegram struct {
	BotToken string
	ChatId   int64
}

type WebhookMessage struct {
	EventMeta EventMeta `json:"eventmeta"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
}

type EventMeta struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Reason    string `json:"reason"`
}

func (t *Telegram) Init(c *config.Config) error {
	botToken := c.Handler.Telegram.BotToken
	if botToken == "" {
		botToken = os.Getenv("KW_TELEGRAM_BOT_TOKEN")
	}
	t.BotToken = botToken

	var err error
	chatId := c.Handler.Telegram.ChatId
	if chatId == 0 && os.Getenv("KW_TELEGRAM_CHAT_ID") != "" {
		chatId, err = strconv.ParseInt(os.Getenv("KW_TELEGRAM_CHAT_ID"), 10, 64)
		if err != nil {
			return err
		}
	}
	t.ChatId = chatId

	return checkMissingVars(t)
}

func (t *Telegram) Handle(e event.Event) {
	botApi, err := tgbotapi.NewBotAPI(t.BotToken)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	_, err = botApi.Send(tgbotapi.NewMessage(t.ChatId, e.Message()))
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to %v at %s ", t.ChatId, time.Now())
}

func checkMissingVars(t *Telegram) error {
	if t.BotToken == "" {
		return fmt.Errorf(telegramErrMsg, "Missing Telegram bot token")
	}
	if t.ChatId == 0 {
		return fmt.Errorf(telegramErrMsg, "Missing Telegram chat id")
	}
	return nil
}
