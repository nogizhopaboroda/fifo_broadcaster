package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	// "encoding/binary"
	"flag"
	// "fmt"
	"log"
	"os"
	// "io"
	"./tcp_broadcaster"
	"./ws_broadcaster"
)

var messages = make(chan []byte)

func Broadcast() {
	for {
		select {
		case message := <-messages:
			wsBroadcaster.Broadcast(message)
			tcpBroadcaster.Broadcast(message)
		}
	}
}

type reader interface {
	ReadString(delim byte) (line string, err error)
}

func read(r reader, delim []byte) (line []byte, err error) {
	for {
		s := ""
		s, err = r.ReadString(delim[len(delim)-1])
		if err != nil {
			return
		}

		line = append(line, []byte(s)...)
		if bytes.HasSuffix(line, delim) {
			return line[:len(line)-len(delim)], nil
		}
	}
}

func readSource(fileName string, bufferSize int, chunkSeparator string, discardSeparator bool) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	//Experimenting with buffers
	//TODO refactor this somehow
	if bufferSize == 0 {
		// separator := []byte{byte(0), byte(0), byte(0), byte(1)}
		separator, err := hex.DecodeString(chunkSeparator)
		if err != nil {
			log.Fatal("Invalid separator: ", err)
		}
		for {
			line, err := read(reader, separator)
			var chunk []byte
			if discardSeparator == true {
				chunk = line
			} else {
				chunk = append(separator, line...)
			}

			if err == nil {
				// fmt.Println(line)
				messages <- chunk
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
	chunkSeparator := flag.String("s", "0a", "chunks separator")
	discardSeparator := flag.Bool("ds", false, "discard separator")
	flag.Parse()

	go wsBroadcaster.Start(*wsPortNo)
	go tcpBroadcaster.Start(*tcpPortNo)

	go Broadcast()

	readSource(*fileName, *bufferSize, *chunkSeparator, *discardSeparator)
}
