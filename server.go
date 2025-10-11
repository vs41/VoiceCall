package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

func main() {
	app := fiber.New()
	app.Static("/", "./public")
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("i3.html", fiber.Map{
			"title": "Team Chat",
		})
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			mutex.Lock()
			delete(clients, c)
			mutex.Unlock()
			log.Println("🔌 Client disconnected")
			c.Close()
		}()

		log.Println("🔗 New client connected")

		mutex.Lock()
		clients[c] = true
		mutex.Unlock()

		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			log.Printf("📨 Message from client: %s\n", msg)

			mutex.Lock()
			for conn := range clients {
				if conn != c {
					if err := conn.WriteMessage(msgType, msg); err != nil {
						log.Println("❌ Write error to client:", err)
						conn.Close()
						delete(clients, conn)
					}
				}
			}
			mutex.Unlock()
		}
	}))

	log.Println("🚀 Server started on ws://0.0.0.0:8001/ws")
	log.Fatal(app.Listen(":8022"))
}
