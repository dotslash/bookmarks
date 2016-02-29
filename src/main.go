package main

import (
    "net/http"
    "os"
    "fmt"
)

var server_prefix string
func main() {
    fmt.Println("starting")
    argsWithoutProg := os.Args[1:]
    Log.Info("args", argsWithoutProg)
    server_prefix = argsWithoutProg[0]
    router := NewRouter()
    Log.Warn(http.ListenAndServe(":8085", router))
}
