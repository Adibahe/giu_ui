package main

import (
	"database/sql"
	"flag"
	"log"

	g "github.com/AllenDang/giu"
)

var db *sql.DB
var recent = make(map[string]message)

var monoFont *g.FontInfo

func main() {

	test := flag.Bool("test", false, "enables testing for ui by sending fake messages")

	flag.Parse()

	msgchan := make(chan message, 100)
	var messages []message

	// if err := g.ParseCSSStyleSheet(cssStyle); err != nil {
	// 	panic(err)
	// }

	db = connectDb()
	defer db.Close()

	if !*test {
		go startServer(msgchan)
	} else {
		log.Println("in testingUi")
		go testingUi(msgchan)
	}

	wnd := g.NewMasterWindow("Function Call Viewer", 800, 600, 0)
	monoFont = g.Context.FontAtlas.AddFont("C:/Windows/Fonts/consola.ttf", 16)

	wnd.Run(func() {
		messages = loop(msgchan, messages)
	})
}
