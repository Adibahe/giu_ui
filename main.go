package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"giu_ui/inject"

	g "github.com/AllenDang/giu"
	webview "github.com/webview/webview_go"
)

var db *sql.DB
var recent = make(map[string]message)
var monoFont *g.FontInfo

func main() {

	test := flag.Bool("test", false, "enables testing for ui by sending fake messages")

	flag.Parse()

	msgchan := make(chan message, 100)
	var messages []message

	db = connectDb()
	defer db.Close()

	go startServer() // hosts webpage on localhost

	// starts webview
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("BinStop")
	w.SetSize(900, 700, webview.HintNone)
	waitForServer("http://localhost:8080")
	w.Navigate(fmt.Sprintf("http://localhost:8080/?t=%d", time.Now().UnixNano()))

	if !*test {
		go openPipe(msgchan) // used to get data from backend (eg. function calls, id)
	} else {
		log.Println("in testingUi")
		go testingUi(msgchan)
	}

	w.Bind("openExternalLink", openExternalLink)
	uiUpdater(w, msgchan, &messages)
	w.Bind("onPageReload", func() {
		log.Println("Webpage Reloaded")
		for _, msg := range messages {
			js := fmt.Sprintf(`window.addRow(%q, %q)`, msg.Id, msg.Name)
			w.Eval(js)
		}
	})
	w.Bind("debugCommand", DebugCommand)
	w.Bind("openFileDialog", OpenFileDialog)
	giveToJs(w)
	w.Bind("injectHookDll", inject.InjectHookDll)

	w.Run()

}
