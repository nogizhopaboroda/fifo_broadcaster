package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"log"
	"os"
)

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

type ReaderParams struct {
	fileName         string
	bufferSize       int
	chunkSeparator   string
	discardSeparator bool
	channel          chan []byte
}

func ReadSource(readerParams ReaderParams) {
	fileName := readerParams.fileName
	bufferSize := readerParams.bufferSize
	chunkSeparator := readerParams.chunkSeparator
	discardSeparator := readerParams.discardSeparator
	channel := readerParams.channel

	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	//Experimenting with buffers
	//TODO refactor this somehow
	if bufferSize == 0 {
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
				channel <- chunk
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
