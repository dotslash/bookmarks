package main

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strings"
    "fmt"
)


func ActionView(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    aliasinfos := getAllAliases()
    resp := makeResponse(aliasinfos, server_prefix)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        Log.Fatal(err)
    }
}

func ActionAdd(w http.ResponseWriter, r *http.Request) {
    short := r.FormValue("short")
    long := r.FormValue("url")
    secret := r.FormValue("secret")
    Log.Println(short, long, secret)
    respStr := addAlias(long, short, secret)
    fmt.Fprint(w, respStr)
}

func ActionDel(w http.ResponseWriter, r *http.Request) {
    short := r.FormValue("id")
    secret := r.FormValue("secret")
    Log.Println("formparams", short, secret)
    respStr := delAlias(short, secret)
    fmt.Fprint(w, respStr)
}

func ActionUpdate(w http.ResponseWriter, r *http.Request) {
    presAlias := r.FormValue("id")
    newVal := r.FormValue("newvalue")
    oldVal := r.FormValue("oldvalue")
    colname := r.FormValue("colname")
    secret := r.FormValue("secret")
    Log.Println("formparams:(presAlias, new, old, field)", presAlias, newVal, oldVal, colname)
    respStr := updateAlias(presAlias, oldVal, newVal, colname, secret)
    fmt.Fprint(w, respStr)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    redId := vars["redId"]
    url := urlFromAlias(redId)
    if url != nil {
        var urlStr = *url
        if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
            urlStr = "http://" + urlStr
        }
        http.Redirect(w, r, urlStr, http.StatusFound)
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
        panic(err)
    }
}
