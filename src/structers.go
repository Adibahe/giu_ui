package main

type message struct {
	Id      string `json: "id"`
	Name    string `json: "name"`
	Content string
}

type response struct {
	Ok  bool   `json: "ok"`
	Msg string `json: "msg"`
}
