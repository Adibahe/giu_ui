package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	winio "github.com/Microsoft/go-winio"
)

func handleConn(conn net.Conn, msgchan chan message) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var msg message
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("decode failed: %v", err)
	}
	fmt.Printf("id: %s funcName: %s \n", msg.Id, msg.Name)
	msgchan <- msg // adding messages to channel

	var resp response
	resp.Ok = true
	resp.Msg = "got it!!"

	err = encoder.Encode(&resp)
	if err != nil {
		log.Fatalf("Encoder failed: %v", err)
	}

}

func startServer(msgchan chan message) {
	log.Println("Started Server!!")

	const pipename = `\\.\pipe\giu_ui_Pipe`

	ln, err := winio.ListenPipe(pipename, nil)
	if err != nil {
		log.Fatalf("Listening failed: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept() failed: %v", err)
			continue
		}

		handleConn(conn, msgchan)

	}
}
