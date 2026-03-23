package main

import (
	"database/sql"

	g "github.com/AllenDang/giu"
)

var db *sql.DB
var recent = make(map[string]message)

func main() {
	msgchan := make(chan message, 100)
	var messages []message

	// if err := g.ParseCSSStyleSheet(cssStyle); err != nil {
	// 	panic(err)
	// }

	db = connectDb()
	defer db.Close()
	go startServer(msgchan)

	wnd := g.NewMasterWindow("Function Call Viewer", 800, 600, 0)
	wnd.Run(func() {
		messages = loop(msgchan, messages)
	})
}
