package main

import (
	"log"
	"net/http"
	"os"
)
var server_prefix string
func main() {
	argsWithoutProg := os.Args[1:]
	log.Println("args", argsWithoutProg)
	server_prefix = argsWithoutProg[0]
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8085", router))
}
