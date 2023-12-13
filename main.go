package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8089", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type Channel struct {
	subscribers map[*websocket.Conn]bool
	mutex       sync.Mutex
}

var channels = make(map[string]*Channel)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer wsConn.Close()

	log.Println("Client Connected")
	reader(wsConn)
}

func reader(conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Handle message for subscribing, unsubscribing, or sending to a channel
		// Example: "subscribe:channel1" or "message:channel1:Hello World"
		handleMessage(conn, mt, message)
	}
}

// subscribe:<channel_name>
// unsubscribe:<channel_name>
// message:<channel_name>:<message_content>
func handleMessage(conn *websocket.Conn, mt int, message []byte) {
	msgParts := strings.Split(string(message), ":")
	if len(msgParts) < 2 {
		return // Invalid message format
	}

	command := msgParts[0]
	channelName := msgParts[1]

	switch command {
	case "subscribe":
		subscribeToChannel(conn, channelName)

	case "unsubscribe":
		unsubscribeFromChannel(conn, channelName)

	case "message":
		if len(msgParts) < 3 {
			return // Invalid message format for sending message
		}
		broadcastMessage(channelName, mt, []byte(strings.Join(msgParts[2:], ":")))
	}
}

func subscribeToChannel(conn *websocket.Conn, channelName string) {
	channelsMutex.Lock()
	defer channelsMutex.Unlock()

	if channels[channelName] == nil {
		channels[channelName] = &Channel{subscribers: make(map[*websocket.Conn]bool)}
	}
	channels[channelName].subscribers[conn] = true
}

func unsubscribeFromChannel(conn *websocket.Conn, channelName string) {
	channelsMutex.Lock()
	defer channelsMutex.Unlock()

	if channel, ok := channels[channelName]; ok {
		delete(channel.subscribers, conn)
	}
}

func broadcastMessage(channelName string, mt int, message []byte) {
	channelsMutex.RLock()
	defer channelsMutex.RUnlock()

	channel, ok := channels[channelName]
	if !ok {
		return // Channel does not exist
	}

	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	for conn := range channel.subscribers {
		if err := conn.WriteMessage(mt, message); err != nil {
			log.Println("write error:", err)
			// Optional: handle failed delivery, e.g., by removing the subscriber
		}
	}
}

// Define channelsMutex as a global variable for thread-safe access to the channels map
var channelsMutex sync.RWMutex

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	setupRoutes()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
