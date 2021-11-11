package main

import (
	"bufio"
	"bytes"
	"errors"
	"go.bug.st/serial"
	"io"
	"strings"
)

const startRpcSessionCommand = "start_rpc_session\r"

func initCli(port string) (io.ReadWriter, error) {
	ser, err := serial.Open(port, &serial.Mode{})
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(ser)
	_, err = readUntil(br, []byte("\r\n\r\n>: "))
	if err != nil {
		return nil, err
	}
	_, err = ser.Write([]byte(startRpcSessionCommand))
	if err != nil {
		return nil, err
	}
	token, err := br.ReadString('\r')
	if err != nil {
		return nil, err
	}
	if token != startRpcSessionCommand {
		return nil, errors.New(strings.TrimSpace(token))
	}
	return ser, nil
}

type reader interface {
	ReadString(delim byte) (line string, err error)
}

func readUntil(r reader, delim []byte) (line []byte, err error) {
	for {
		s := ""
		s, err = r.ReadString(delim[len(delim)-1])
		if err != nil {
			return
		}

		line = append(line, []byte(s)...)
		if bytes.HasSuffix([]byte(s), delim) {
			return line[:len(line)-len(delim)], nil
		}
	}
}
