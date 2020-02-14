package tcpBroadcaster

import (
    "log"
    "fmt"
    "net"
)


type hub struct {
	messages    chan []byte
	connections map[int]net.Conn
	addConn     chan net.Conn
}

var h = &hub{
	messages:    make(chan []byte),
	connections: make(map[int]net.Conn),
	addConn:     make(chan net.Conn),
}

func (h *hub) sendAll(message []byte) {
	// expired := []int{}
	for _, conn := range h.connections {
		conn.Write([]byte(message))
		// fmt.Println("data")
		// if err != nil {
			// conn.Close()
			// expired = append(expired, i)
		// }
	}
	// // Prune the obsolete connections
	// if len(expired) > 0 {
		// for _, connId := range expired {
			// log.Println("Closed connection:", connId)
			// delete(h.connections, connId)
		// }
	// }
}

func Broadcast(message []byte) {
  h.messages <- []byte(message)
}

func (h *hub) Run() {
        counter := 0

        for {
            select {
            case conn := <-h.addConn:
                h.connections[counter] = conn
                counter = counter + 1
                log.Println("TCP client connected:", counter)
            case c := <-h.messages:
		h.sendAll([]byte(c))
            }
        }
}

func connectionHandler(ln net.Listener) {
        conn, err := ln.Accept()
        if err != nil {
            panic(err)
        }

	h.addConn <- conn
}

func Serve(port string, connectionHandler func (net.Listener)) {
    ln, err := net.Listen("tcp", port)
    if err != nil {
        panic(err)
    }


    // connection queue
    for {
      connectionHandler(ln)
    }
}


func Start(portNo int){
  go h.Run()
  log.Printf("Serving TCP on %d...", portNo)
  Serve(fmt.Sprintf(":%d", portNo), connectionHandler)
}
