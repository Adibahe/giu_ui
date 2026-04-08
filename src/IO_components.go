package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func testingUi(msgchan chan message) {
	log.Println("in testingUi")
	ids := []string{"604", "166", "8", "868", "869", "12", "189", "461", "148", "256", "205", "153", "743"}
	i := 0
	for _, id := range ids {

		var msg message
		// n := rand.Int64N(1000)
		msg.Id = id
		msg.Name = ""

		msgchan <- msg

		time.Sleep(time.Second)
		i++
	}
}

func handleConn(conn net.Conn, msgchan chan message) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	for {
		var msg message
		err := decoder.Decode(&msg)
		if err != nil {
			log.Printf("decode failed: %v", err)
			break
		}
		fmt.Printf("id: %s funcName: %s \n", msg.Id, msg.Name)
		msgchan <- msg
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

func waitForServer(url string) {
	for i := 0; i < 20; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Fatal("server did not start in time")

}
