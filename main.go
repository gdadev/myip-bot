package main

import (
	"log"
	"net"
	"time"

	"github.com/kkyr/fig"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Config struct {
	Token       string `validate:"required"`
	IDWhitelist []int64
}

func main() {
	var cfg Config

	if err := fig.Load(&cfg, fig.UseEnv("MYIP_BOT"), fig.IgnoreFile()); err != nil {
		log.Fatalf("Parsing config: %v", err)
	}

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(e error, c tele.Context) {
			log.Println(e)
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Creating bot: %v", err)
	}

	var (
		menu    = &tele.ReplyMarkup{ResizeKeyboard: true}
		btnMyIP = menu.Text("My IP")
	)

	if len(cfg.IDWhitelist) > 0 {
		b.Use(middleware.Whitelist(cfg.IDWhitelist...))
	}

	menu.Reply(
		menu.Row(btnMyIP),
	)

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Welcome", menu)
	})

	b.Handle(&btnMyIP, func(c tele.Context) error {
		conn, err := net.Dial("udp", "1.1.1.1:80")
		if err != nil {
			return err
		}

		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return c.Send(localAddr.IP.To4().String())
	})

	b.Start()
}
