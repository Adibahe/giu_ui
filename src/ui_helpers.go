package main

import (
	"fmt"
	"log"

	webview "github.com/webview/webview_go"
)

func addTableRow(w webview.WebView, id string, functionName string) {
	functionName = getName(id)
	w.Dispatch(func() {
		js := fmt.Sprintf(`window.addRow(%q, %q)`, id, functionName)
		w.Eval(js)
		fmt.Printf("sent: %s, %s \n", id, functionName)
	})
}

func uiUpdater(w webview.WebView, msgChan <-chan message) {
	go func() {
		for msg := range msgChan {
			addTableRow(w, msg.Id, msg.Name)
		}
	}()
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
