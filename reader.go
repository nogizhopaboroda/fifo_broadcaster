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
	fileName       string
	bufferSize     int
	chunkSeparator string
	keepSeparator  string
	channel        chan []byte
}

func ReadSource(readerParams ReaderParams) {
	fileName := readerParams.fileName
	bufferSize := readerParams.bufferSize
	chunkSeparator := readerParams.chunkSeparator
	keepSeparator := readerParams.keepSeparator
	channel := readerParams.channel

	var file *os.File

	if fileName == "-" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.OpenFile(fileName, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatal("Open named pipe file error:", err)
		}

		defer file.Close()
	}

	reader := bufio.NewReader(file)

	if bufferSize == 0 {
		separator, err := hex.DecodeString(chunkSeparator)
		if err != nil {
			log.Fatal("Invalid separator: ", err)
		}

		var buff []byte
		for {
			line, err := read(reader, separator)
			var chunk []byte
			switch keepSeparator {
			case "none":
				chunk = line
			case "beginning-of-next":
				chunk = append(buff, line...)
				buff = append(separator[:0:0], separator...)
			case "end-of-current":
				chunk = append(line, separator...)
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
