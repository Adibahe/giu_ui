package main

import (
	g "github.com/AllenDang/giu"
)

func main() {
	msgchan := make(chan message, 100)
	var messages []message

	wnd := g.NewMasterWindow("Function Call Viewer", 800, 600, 0)

	// if err := g.ParseCSSStyleSheet(cssStyle); err != nil {
	// 	panic(err)
	// }

	go startServer(msgchan)

	wnd.Run(func() {
		messages = loop(msgchan, messages)
	})
}
