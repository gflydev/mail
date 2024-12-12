package main

import (
	"github.com/gflydev/core"
	"github.com/gflydev/mail"
)

// =========================================================================================
//                                     Application
// =========================================================================================

func main() {
	app := core.New()

	mail.Send(mail.Envelop{
		To:      []string{"vohuynhvinh@gmail.com"},
		ReplyTo: []string{"admin@jivecode.com"},
		Subject: "Test mail from gflydev/mail",
		Text:    "Test mail from gflydev/mail",
		HTML:    "<h2>Test mail from gflydev/mail</h2>",
	})

	app.Run()
}
