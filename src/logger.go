package main

import (
    "net/http"
    "time"
    log "github.com/Sirupsen/logrus"
    "github.com/rifflock/lfshook"
    "os/user"
)
var Log *log.Logger = newLogger()

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        inner.ServeHTTP(w, r)
        Log.Infof(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.RequestURI,
            name,
            time.Since(start),
        )
    })
}



func newLogger() *log.Logger {
    usr, _ := user.Current()
    home := usr.HomeDir
    lLog := log.New()
    lLog.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
        log.InfoLevel : home + "/log/bm-info.log",
        log.ErrorLevel : home + "/log/bm-error.log",
    }))
    return lLog
}
