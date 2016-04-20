package main

import (
	"fmt"
	"net/http"
	"os"
)

var server_prefix string

func main() {
	fmt.Println("starting")
	argsWithoutProg := os.Args[1:]
	Log.Info("args", argsWithoutProg)
	server_prefix = argsWithoutProg[0]
	port := argsWithoutProg[1]
	router := NewRouter()
	Log.Warn(http.ListenAndServe(":"+port, router))
}
