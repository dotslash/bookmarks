package main

import (
	"fmt"
	"net/http"
	"os"
)

// ServerAddress is server address the bookmarks server is running as.
// TODO(dotslash): Inject this into the router as opposed to making
// this a constant.
var ServerAddress string

func main() {
	fmt.Println("starting")
	argsWithoutProg := os.Args[1:]
	Log.Info("args", argsWithoutProg)
	ServerAddress = argsWithoutProg[0]
	port := argsWithoutProg[1]
	router := NewRouter()
	Log.Warn(http.ListenAndServe(":"+port, router))
}
