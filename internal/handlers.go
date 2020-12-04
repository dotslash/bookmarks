package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/vincent-petithory/dataurl"
)

// Handlers struct has methods to handle http requests to be served by
// the webserver.
type Handlers struct {
	once    sync.Once
	storage *StorageInterface
	// Address of the server hosting the application. This will be used
	// to generate short urls. E.g - https://bm.suram.in/r/foo
	serverAddress string
	// Location of the sqlite file.
	dbFile string
}

func (h *Handlers) getStorage() *StorageInterface {
	h.once.Do(func() {
		h.storage = NewStorageInterface(h.dbFile)
	})
	return h.storage
}

// ActionView handles http request to view bookmarks list.
func (h *Handlers) ActionView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	secret := r.FormValue("secret")
	aliasInfos := h.getStorage().GetAllAliases(secret)
	resp := CreateViewResponse(aliasInfos, h.serverAddress)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		Log.Fatal(err)
	}
}

// ActionLookup handles http request to convert short url to the full url.
func (h *Handlers) ActionLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	short := r.FormValue("short")
	fullURL := h.getStorage().URLFromAlias(short)
	if fullURL == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
		return
	}
	fmt.Fprint(w, *fullURL)
}

// ActionRevLookup handles http request to convert full url to the short url.
func (h *Handlers) ActionRevLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	secret := r.FormValue("secret")
	long := r.FormValue("long")
	shortUrls := h.getStorage().GetShortUrls(secret, long)
	resp := RevLookupResponse{shortUrls}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		Log.Fatal(err)
	}
}

// ActionAdd handles http request to add a new bookmark.
func (h *Handlers) ActionAdd(w http.ResponseWriter, r *http.Request) {
	short := r.FormValue("short")
	long := r.FormValue("url")
	secret := r.FormValue("secret")
	Log.Println(short, long, secret)
	respStr := h.getStorage().AddAlias(long, short, secret)
	fmt.Fprint(w, respStr)
}

// ActionDel handles http request to add a delete bookmark.
func (h *Handlers) ActionDel(w http.ResponseWriter, r *http.Request) {
	short := r.FormValue("id")
	secret := r.FormValue("secret")
	Log.Println("formparams", short, secret)
	respStr := h.getStorage().DelAlias(short, secret)
	fmt.Fprint(w, respStr)
}

// ActionUpdate handles http request to add a update bookmark.
func (h *Handlers) ActionUpdate(w http.ResponseWriter, r *http.Request) {
	presAlias := r.FormValue("id")
	newVal := r.FormValue("newvalue")
	oldVal := r.FormValue("oldvalue")
	colName := r.FormValue("colname")
	secret := r.FormValue("secret")
	Log.Println("formparams:(presAlias, new, old, field)", presAlias, newVal, oldVal, colName)
	respStr := h.getStorage().UpdateAlias(presAlias, oldVal, newVal, colName, secret)
	fmt.Fprint(w, respStr)
}


func (h *Handlers) writeError(w http.ResponseWriter, status int, text string ) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(ErrStruct{Code: status, Text: text})
	if err != nil {
		panic(err)
	}
}

func (h *Handlers) handleDataUrl(w http.ResponseWriter, urlStr string) {
	parsed, err := dataurl.DecodeString(urlStr)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Add("Content-Type", parsed.MediaType.ContentType())
	w.Write(parsed.Data)
}


// Redirect handles redirect for short url to full url.
func (h *Handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	redID := vars["redId"]
	url := h.getStorage().URLFromAlias(redID)
	if url != nil {
		var urlStr = *url
		if strings.HasPrefix(urlStr, "data:") {
			h.handleDataUrl(w, urlStr)
			return
		}
		if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
			urlStr = "http://" + urlStr
		}
		http.Redirect(w, r, urlStr, http.StatusFound)
		return
	}
	// If we didn't find it, 404
	h.writeError(w, http.StatusNotFound, "Not Found")
}
