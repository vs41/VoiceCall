// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"
	"webrtc/sever"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/pion/logging"
)

var (
	indexTemplate = &template.Template{}
	log           = logging.NewDefaultLoggerFactory().NewLogger("sfu-fiber")
)

func main() {
	// Read index.html
	indexHTML, err := os.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	indexTemplate = template.Must(template.New("").Parse(string(indexHTML)))

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://v222-j7vv.onrender.com",
	}))
	// Serve index.html
	app.Get("/", func(c *fiber.Ctx) error {
		// gameID := c.Query("gameID")
		// teamID := c.Query("teamID")
		// username := c.Query("username")

		wsURL := "wss://v222-j7vv.onrender.com/websocket?gameID=28&teamID=1&unitID=47&encyptionOn=0&username=blue1&radioRange=40"
		fmt.Println(wsURL)
		var buf bytes.Buffer
		if err := indexTemplate.Execute(&buf, wsURL); err != nil {
			return err
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.Send(buf.Bytes())
	})

	app.Use("/websocket", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket handler
	app.Get("/websocket", websocket.New(sever.WebsocketHandler))

	// Periodically request keyframes
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			sever.DispatchAllKeyFrames()
		}
	}()

	// tlsCert := "cert.pem"
	// tlsKey := "key.pem"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// tlsCert := "/home/vishal_s/Documents/key/cert.pem"
	// tlsKey := "/home/vishal_s/Documents/key/key.pem"
	if err := app.ListenTLS(":"+port); err != nil {
		log.Errorf("Failed to start Fiber server: %v", err)
	}
	// if err := app.ListenTLS(":"+port, tlsCert, tlsKey); err != nil {
	// 	log.Errorf("Failed to start Fiber server: %v", err)
	// }
}
