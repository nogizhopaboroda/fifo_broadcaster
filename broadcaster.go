package main

import (
	"bufio"
	// "encoding/binary"
	"flag"
	"log"
	"os"
	// "fmt"
	// "io"
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


func readSource(fileName string, bufferSize int) {



	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)


	//Experimenting with buffers
	//TODO refactor this somehow
	if bufferSize == 0 {
	  for {
	    line, err := reader.ReadBytes('\n')
	    if err == nil {
	      messages <- []byte(line)
	    }
	  }
	} else {
	  buffer := make([]byte, bufferSize)

	  for {
	    bytesread, err := reader.Read(buffer)

	    if err == nil {
	      messages <- []byte(buffer[0:bytesread])
	    }
	  }
	}

}


func main() {
  bufferSize := flag.Int("b", 0, "buffer size (in bytes)")
  wsPortNo := flag.Int("wp", 8080, "websocket server port")
  tcpPortNo := flag.Int("tp", 1235, "tcp server port")
  fileName := flag.String("i", "", "input source")
  flag.Parse()

  go wsBroadcaster.Start(*wsPortNo)
  go tcpBroadcaster.Start(*tcpPortNo)

  go Broadcast()

  readSource(*fileName, *bufferSize)
}
