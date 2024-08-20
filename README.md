# gFly Mail - Send mail

    Copyright © 2023, gFly
    https://www.gfly.dev
    All rights reserved.

### Usage

Install
```bash
go get -u github.com/gflydev/mail@v1.0.0
```

Quick usage `main.go`
```go
import (
    "github.com/gflydev/mail"	
)

func main() {
    mail.Send(mail.Envelop{
        To:      []string{"vinh@gfly.dev"},
        ReplyTo: []string{"vinh@jivecode.com"},
        Subject: "Test mail from gflydev",
        Text:    "Test mail from gflydev",
        HTML:    "<h2>Test mail from gflydev</h2>",
    })
}
```
