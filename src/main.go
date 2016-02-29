package main

import (
    "log"
    "net/http"
    "os"
)
var server_prefix string
func main() {
    argsWithoutProg := os.Args[1:]
    f, _ := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    defer f.Close()
    log.SetOutput(f)
    log.Println("args", argsWithoutProg)
    server_prefix = argsWithoutProg[0]
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8085", router))
}
