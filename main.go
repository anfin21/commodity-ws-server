package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"log"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"commodity-ws-server/configs"
	"commodity-ws-server/pkg/logger"
)

// Client represents a connected client.
type Client struct {
	conn     net.Conn
	channels map[string]struct{}
	mu       sync.Mutex
	log      *logger.Logger
}

// NewClient creates a new Client.
func NewClient(conn net.Conn, log *logger.Logger) *Client {
	return &Client{
		conn:     conn,
		channels: make(map[string]struct{}),
		log:      log,
	}
}

// Subscribe adds a channel to the client's subscription list.
func (c *Client) Subscribe(channel string) {
	c.log.Debugf("Client begin subscribe channel: %v", channel)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.channels[channel] = struct{}{}
	c.log.Debugf("Client subscribed channel: %v", channel)
}

// Unsubscribe removes a channel from the client's subscription list.
func (c *Client) Unsubscribe(channel string) {
	c.log.Debugf("Client begin unsubscribe channel: %v", channel)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.channels, channel)
	c.log.Debugf("Client unsubscribed channel: %v", channel)
}

// Broadcast sends a message to all clients subscribed to a channel.
func Broadcast(clients map[*Client]struct{}, mu *sync.Mutex, channel, message string) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		client.mu.Lock()
		_, subscribed := client.channels[channel]
		client.mu.Unlock()
		if subscribed {
			_ = wsutil.WriteServerMessage(client.conn, ws.OpText, []byte(fmt.Sprintf("%s", message)))
		}
	}
}

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	log := logger.New(cfg)

	clients := make(map[*Client]struct{})
	mu := sync.Mutex{}

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Errorf("WebSocket upgrade error: %v", err)
			return
		}

		client := NewClient(conn, log)

		mu.Lock()
		clients[client] = struct{}{}
		mu.Unlock()

		go handleClient(client, clients, &mu, log)
	})

	http.ListenAndServe(fmt.Sprintf(":%v", cfg.Server.Port), mux)
}

func handleClient(client *Client, clients map[*Client]struct{}, mu *sync.Mutex, log *logger.Logger) {
	defer func() {
		mu.Lock()
		delete(clients, client)
		mu.Unlock()

		client.conn.Close()
	}()

	for {
		msg, _, err := wsutil.ReadClientData(client.conn)
		if err != nil {
			return
		}

		payload := string(msg)
		log.Debugf("Received payload: %s", payload)

		if len(payload) > 10 && payload[:10] == "subscribe:" {
			log.Debugf("Handle payload CASE 1: subscribe")
			channel := payload[10:]
			client.Subscribe(channel)
		} else if len(payload) > 12 && payload[:12] == "unsubscribe:" {
			log.Debugf("Handle payload CASE 2: unsubscribe")
			channel := payload[12:]
			client.Unsubscribe(channel)
		} else if len(payload) > 8 && payload[:8] == "message:" {
			log.Debugf("Handle payload CASE 3: broadcast message")
			parts := strings.SplitN(payload[8:], ":", 2)
			if len(parts) == 2 {
				channel, message := parts[0], parts[1]
				log.Debugf("Handle payload CASE 3: BEGIN broadcast message=%v to channel=%v", message, channel)
				Broadcast(clients, mu, channel, message)
			}
		} else {
			log.Debugf("Handle payload CASE 4: payload does not follow any defined format | payload=%v", payload)
		}
	}
}
