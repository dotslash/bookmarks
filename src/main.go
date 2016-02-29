package main

import (
    "net/http"
    "os"
)

var server_prefix string
func main() {
    argsWithoutProg := os.Args[1:]
    Log.Info("args", argsWithoutProg)
    server_prefix = argsWithoutProg[0]
    router := NewRouter()
    Log.Warn(http.ListenAndServe(":8085", router))
}
