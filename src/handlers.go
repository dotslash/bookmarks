package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func ActionView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	secret := r.FormValue("secret")
	aliasInfos := getAllAliases(secret)
	resp := makeViewResponse(aliasInfos, server_prefix)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		Log.Fatal(err)
	}
}

func ActionLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	short := r.FormValue("short")
	fullUrl := urlFromAlias(short)
	if fullUrl == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
		return
	}
	fmt.Fprint(w, *fullUrl)
}

func ActionRevLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	secret := r.FormValue("secret")
	long := r.FormValue("long")
	shortUrls := getShortUrls(secret, long)
	resp := makeRevLookUpResponse(shortUrls)
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
	colName := r.FormValue("colname")
	secret := r.FormValue("secret")
	Log.Println("formparams:(presAlias, new, old, field)", presAlias, newVal, oldVal, colName)
	respStr := updateAlias(presAlias, oldVal, newVal, colName, secret)
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
