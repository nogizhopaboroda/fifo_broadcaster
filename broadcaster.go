package main

import (
	"bufio"
	"flag"
	"log"
	"os"
        "./ws_broadcaster"
        "./tcp_broadcaster"
)

var messages = make(chan []byte)

func Broadcast(){
  for {
          select {
          case message := <-messages:
              wsBroadcaster.Broadcast(message)
              tcpBroadcaster.Broadcast(message)
          }
  }
}


func readSource(fileName string) {

	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if err == nil {
                  messages <- []byte(line)
		}
	}
}


func main() {
  wsPortNo := flag.Int("wp", 8080, "websocket server port")
  tcpPortNo := flag.Int("tp", 1235, "tcp server port")
  fileName := flag.String("i", "", "input source")
  flag.Parse()

  go wsBroadcaster.Start(*wsPortNo)
  go tcpBroadcaster.Start(*tcpPortNo)

  go Broadcast()

  readSource(*fileName)
}
