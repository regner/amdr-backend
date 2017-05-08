package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type InjestPostRequest struct {
	MarketItems []string `json:"marketItems"`
}

type MarketUpdate struct {
	MarketItems []MarketItem `json:"MarketItems"`
	IngestTime  time.Time    `json:"IngestTime"`
}

type MarketItem struct {
	ID               int    `json:"Id"`
	UnitPrice        int    `json:"UnitPriceSilver"`
	TotalPrice       int    `json:"TotalPriceSilver"`
	Amount           int    `json:"Amount"`
	Tier             int    `json:"Teir"`
	ItemTypeID       string `json:"ItemTypeId"`
	ItemGroupTypeID  string `json:"ItemGroupTypeId"`
	EnchantmentLevel int    `json:"EnchantmentLevel"`
	QualityLevel     int    `json:"QualityLevel"`
	Expires          string `json:"Expires"`
}

func main() {
	registry := NewRegistry()
	go registry.run()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, _ := nats.Connect(natsURL)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	natsTopic := os.Getenv("NATS_TOPIC")
	if natsTopic == "" {
		natsTopic = "amdr-market-updates"
	}

	ec.Subscribe(natsTopic, func(marketUpdate string) {
		registry.broadcast <- []byte(marketUpdate)
	})

	r := gin.Default()

	r.GET("/api/v1/announce/", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := &Client{registry: registry, conn: conn, send: make(chan []byte, 256)}
		client.registry.register <- client

		go client.run()
	})

	r.POST("/api/v1/ingest/", func(c *gin.Context) {
		var incomingRequest InjestPostRequest
		c.BindJSON(&incomingRequest)

		var marketUpdate MarketUpdate
		for _, v := range incomingRequest.MarketItems {
			var item MarketItem
			json.Unmarshal([]byte(v), &item)

			marketUpdate.MarketItems = append(marketUpdate.MarketItems, item)
		}

		marketUpdate.IngestTime = time.Now().UTC()

		ec.Publish(natsTopic, marketUpdate)
	})

	r.Run(":8080")
}
