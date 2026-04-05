package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func testingUi(msgchan chan message) {
	log.Println("in testingUi")

	for {

		var msg message
		n := rand.Int64N(1000)
		msg.Id = strconv.FormatInt(n, 10)
		msg.Name = ""

		msgchan <- msg

		time.Sleep(time.Second)
	}
}

func handleConn(conn net.Conn, msgchan chan message) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg message
		err := decoder.Decode(&msg)
		if err != nil {
			log.Printf("decode failed: %v", err)
		}
		fmt.Printf("id: %s funcName: %s \n", msg.Id, msg.Name)
		msgchan <- msg

		var resp response
		resp.Ok = true
		resp.Msg = "got it!!"

		err = encoder.Encode(&resp)
		if err != nil {
			log.Fatalf("Encoder failed: %v", err)
		}
	}

}

func openPipe(msgchan chan message) {
	log.Println("Started Server!!")

	const pipename = `\\.\pipe\giu_ui_Pipe`

	ln, err := winio.ListenPipe(pipename, nil)
	if err != nil {
		log.Fatalf("Listening failed: %v", err)
	}
	defer ln.Close()

	//	for {
	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("Accept() failed: %v", err)
	}

	handleConn(conn, msgchan)

	//	}
}

func startServer() {
	fs := http.FileServer(http.Dir("./webui"))
	http.Handle("/", fs)

	log.Println("Serving at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
