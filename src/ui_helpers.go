package main

import (
	"fmt"
	"log"
	"os/exec"

	webview "github.com/webview/webview_go"
)

func addTableRow(w webview.WebView, id string, functionName string) {
	w.Dispatch(func() {
		js := fmt.Sprintf(`window.addRow(%q, %q)`, id, functionName)
		w.Eval(js)
		log.Printf("sent: %s, %s \n", id, functionName)
	})
}

func uiUpdater(w webview.WebView, msgChan <-chan message, messages *[]message) {
	go func() {
		for msg := range msgChan {

			msg.Name = getName(msg.Id)

			*messages = append(*messages, msg)
			addTableRow(w, msg.Id, msg.Name)
		}
	}()
}

func onPageReload(messages []message) {
	fmt.Println("Webpage loaded/reloaded")
}

func giveToJs(w webview.WebView) {
	err := w.Bind("getDescFromGo", func(id string) string {
		var details string
		getDesc(id, &details)
		return details
	})
	if err != nil {
		log.Fatal(err)
	}
}

func openExternalLink(url string) string {
	// log.Println("Opening ", url)
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	if err != nil {
		return err.Error()
	}
	return "ok"
}
