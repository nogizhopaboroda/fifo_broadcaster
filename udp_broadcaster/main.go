package udpBroadcaster

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

func connectionHandler(conn *net.UDPConn) {
}

func handleUDPConnection(conn *net.UDPConn) {

// here is where you want to do stuff like read or write to client

buffer := make([]byte, 1024)

n, addr, err := conn.ReadFromUDP(buffer)

fmt.Println("UDP client : ", addr)
fmt.Println("Received from UDP client :  ", string(buffer[:n]))

if err != nil {
log.Fatal(err)
}

// NOTE : Need to specify client address in WriteToUDP() function
//        otherwise, you will get this error message
//        write udp : write: destination address required if you use Write() function instead of WriteToUDP()

// write message back to client
message := []byte("Hello UDP client!")
_, err = conn.WriteToUDP(message, addr)

if err != nil {
log.Println(err)
}

}

func Serve(port string) {

	hostName := "localhost"
portNum := "6000"
service := hostName + ":" + portNum

udpAddr, err := net.ResolveUDPAddr("udp", service)
udpAddr2, err := net.ResolveUDPAddr("udp", hostName + ":" + "3000")

    ln, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
      panic(err)
    }

fmt.Println("UDP server up and listening on port 6000")

    defer ln.Close()


    connection, err := net.DialUDP("udp", udpAddr2, udpAddr)

    for {
	select {
	case c := <-h.messages:
	    connection.Write([]byte(c))
	    // fmt.Println(c)
	    // h.sendAll([]byte(c))
	}
    }

    // connection queue
    // for {
      // handleUDPConnection(ln)
      // // connectionHandler(ln)
    // }
}


func Start(portNo int){
  // go h.Run()
  log.Printf("Serving UDP on %d...", portNo)
  go Serve(fmt.Sprintf(":%d", portNo))
}
