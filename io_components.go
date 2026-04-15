package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/YUSHACOD/rad_api"

	winio "github.com/Microsoft/go-winio"
)

func testingUi(msgchan chan message) {
	log.Println("in testingUi")
	ids := []string{"604", "166", "8", "868", "869", "12", "189", "461", "148", "256", "205", "153", "743"}
	i := 0
	for _, id := range ids {

		var msg message

		msg.Id = id
		msg.Name = ""

		msgchan <- msg

		time.Sleep(time.Second / 2)
		i++
	}
}

func handleConn(conn net.Conn, msgchan chan message) {
	defer conn.Close()

	for {
		var rawID uint64
		err := binary.Read(conn, binary.LittleEndian, &rawID)
		if err != nil {
			if err != io.EOF {
				log.Printf("failed to read binary ID: %v", err)
			}
			break
		}

		if state.auto_running {
			state.rad.SendCommand(rad_api.CMD_RUN, "")
		}

		var msg message
		msg.Id = strconv.FormatUint(rawID, 10)

		fmt.Printf("id: %s funcName: %s \n", msg.Id, msg.Name)
		msgchan <- msg
	}
}

func openPipe(msgchan chan message) {
	log.Println("Started Server!!")

	const pipename = `\\.\pipe\P7_HOOKS`

	ln, err := winio.ListenPipe(pipename, nil)
	if err != nil {
		log.Fatalf("Listening failed: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Accept() failed: %v", err)
		}

		go handleConn(conn, msgchan)

	}
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
	for range 20 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Fatal("server did not start in time")

}
