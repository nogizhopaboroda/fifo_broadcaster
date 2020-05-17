package wsBroadcaster

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type hub struct {
	messages    chan []byte
	connections map[int]*websocket.Conn
	addConn     chan *websocket.Conn
}

// Initialize our main hub struct
var h = &hub{
	messages:    make(chan []byte),
	connections: make(map[int]*websocket.Conn),
	addConn:     make(chan *websocket.Conn),
}

func (h *hub) sendAll(message []byte) {
	expired := []int{}
	for i, conn := range h.connections {
		err := conn.WriteMessage(websocket.BinaryMessage, []byte(message))
		if err != nil {
			conn.Close()
			expired = append(expired, i)
		}
	}
	// Prune the obsolete connections
	if len(expired) > 0 {
		for _, connId := range expired {
			log.Println("Closed connection:", connId)
			delete(h.connections, connId)
		}
	}
}

func Broadcast(message []byte) {
	h.messages <- []byte(message)
}

func (h *hub) Run() {
	id := 0
	for {
		select {
		// Client has connected
		case c := <-h.addConn:
			log.Println("WS client connected:", id)
			h.connections[id] = c
			id++
		// A new message has been received
		case c := <-h.messages:
			h.sendAll([]byte(c))
		}
	}
}

// Handle upgrades to websocket
func connectionHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	h.addConn <- c
}

func Start(portNo int, tls bool, certFile string, keyFile string) {
	go h.Run()
	http.HandleFunc("/", connectionHandler)
	var err error
	if tls == true {
		log.Printf("Serving secure WS on %d...", portNo)
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", portNo), certFile, keyFile, nil)
	} else {
		log.Printf("Serving WS on %d...", portNo)
		err = http.ListenAndServe(fmt.Sprintf(":%d", portNo), nil)
	}
	log.Fatal(err)
}
