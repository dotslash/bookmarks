package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ActionView handles http request to view bookmarks list.
func ActionView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	secret := r.FormValue("secret")
	aliasInfos := getAllAliases(secret)
	resp := makeViewResponse(aliasInfos, ServerAddress)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		Log.Fatal(err)
	}
}

// ActionLookup handles http request to convert short url to the full url.
func ActionLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	short := r.FormValue("short")
	fullURL := urlFromAlias(short)
	if fullURL == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
		return
	}
	fmt.Fprint(w, *fullURL)
}

// ActionRevLookup handles http request to convert full url to the short url.
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

// ActionAdd handles http request to add a new bookmark.
func ActionAdd(w http.ResponseWriter, r *http.Request) {
	short := r.FormValue("short")
	long := r.FormValue("url")
	secret := r.FormValue("secret")
	Log.Println(short, long, secret)
	respStr := addAlias(long, short, secret)
	fmt.Fprint(w, respStr)
}

// ActionDel handles http request to add a delete bookmark.
func ActionDel(w http.ResponseWriter, r *http.Request) {
	short := r.FormValue("id")
	secret := r.FormValue("secret")
	Log.Println("formparams", short, secret)
	respStr := delAlias(short, secret)
	fmt.Fprint(w, respStr)
}

// ActionUpdate handles http request to add a update bookmark.
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

// Redirect handles redirect for short url to full url.
func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	redID := vars["redId"]
	url := urlFromAlias(redID)
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
