package main

import (
	"./tcp_broadcaster"
	"./ws_broadcaster"
	// "flag"
	"github.com/jessevdk/go-flags"
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

var opts struct {
	InputSource string `short:"i" long:"input" description:"Input source" required:"true"`

	BufferSize int `short:"b" long:"buffer-size" description:"Buffer size" default:"0"`

	WsPortNumber int `short:"w" long:"ws-port" description:"WebSocker server port number" default:"8080"`

	TcpPortNumber int `short:"t" long:"tcp-port" description:"TCP server port number" default:"1235"`

	ChunkSeparator string `short:"s" long:"chunk-separator" description:"Chunk separator"`

	DiscardSeparator bool `short:"d" long:"discard-separator" description:"Discard separator"`
}

func main() {

	_, err := flags.Parse(&opts)

	if err != nil {
		panic(err)
	}

	fileName := opts.InputSource

	bufferSize := opts.BufferSize
	wsPortNo := opts.WsPortNumber
	tcpPortNo := opts.TcpPortNumber
	chunkSeparator := opts.ChunkSeparator
	discardSeparator := opts.DiscardSeparator

	// flag.Parse()

	go wsBroadcaster.Start(wsPortNo)
	go tcpBroadcaster.Start(tcpPortNo)

	go Broadcast()

	params := ReaderParams{
		fileName:         fileName,
		bufferSize:       bufferSize,
		chunkSeparator:   chunkSeparator,
		discardSeparator: discardSeparator,
		channel:          messages,
	}

	ReadSource(params)
}
