// Entry point to the bookmarks server.
// Example usage: ./bookmarks.bin https://bm.suram.in 8085
// arg1 : Address of the server hosting the application. This will be used
//        to generate short urls. E.g - https://bm.suram.in/r/foo
// arg2 : port to run the http server at.
package main

import (
	"fmt"
	"github.com/dotslash/bookmarks/internal"
	"net/http"
	"os"
)

func main() {
	fmt.Println("starting")
	argsWithoutProg := os.Args[1:]
	internal.Log.Info("args", argsWithoutProg)
	// Get port and server address.
	ServerAddress := argsWithoutProg[0]
	port := argsWithoutProg[1]
	// Get dbFile path.
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("os.Getwd failed. %v", err))
	}
	dbFile := pwd + "/foo.db"
	// Launch server.
	router := internal.NewRouter(ServerAddress, dbFile)
	internal.Log.Warn(http.ListenAndServe(":"+port, router))
}
