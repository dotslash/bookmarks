package internal

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// routeStruct contains information about the route and
// also the route's http handler/
type routeStruct struct {
	// Name of the route.
	Name string
	// HTTP methods supported by the route.
	// E.g - GET, POST etc.
	Methods []string
	// Url Pattern that this route can handle.
	Pattern string
	// http Handler for the route.
	HandlerFunc http.HandlerFunc
}

func getRoutes(serverAddress string, dbFile string) []routeStruct {
	getOnly := []string{"GET"}
	getAndPost := []string{"GET", "POST"}
	handlers := &Handlers{serverAddress: serverAddress, dbFile: dbFile, storage: NewStorageInterface(dbFile)}
	return []routeStruct{
		{"ActionAdd", getAndPost, "/actions/add", handlers.ActionAdd},
		{"ActionDel", getAndPost, "/actions/delete", handlers.ActionDel},
		{"ActionLookup", getAndPost, "/actions/lookup", handlers.ActionLookup},
		{"ActionRevLookup", getAndPost, "/actions/revlookup", handlers.ActionRevLookup},
		{"ActionUpdate", getAndPost, "/actions/update", handlers.ActionUpdate},
		{"ActionView", getAndPost, "/actions/view", handlers.ActionView},
		{"Redirect", getOnly, "/r/{redId:.*}", handlers.Redirect},
	}
}

// NewRouter returns a new mux.Router that handles all incoming http requests.
// serverAddress: The server address the bookmarks server is running as.
//                E.g - https://bm.suram.in
func NewRouter(serverAddress string, dbFile string) *mux.Router {
	router := mux.NewRouter()
	for _, route := range getRoutes(serverAddress, dbFile) {
		var handler http.Handler
		Log.Println(route)
		handler = route.HandlerFunc
		handler = HTTPLogger(handler, route.Name)
		router.Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	isDev := strings.HasPrefix(serverAddress, "http://localhost") || strings.HasPrefix(serverAddress, "http://127.0.0.1")
	router.PathPrefix("/").Handler(CacheMiddleware(http.FileServer(http.Dir("./internal/static")), isDev))
	return router
}

// CacheMiddleware sets Cache-Control headers for static assets based on environment.
func CacheMiddleware(h http.Handler, isDev bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDev {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		} else {
			// Cache for 30 days in prod
			w.Header().Set("Cache-Control", "public, max-age=2592000")
		}
		h.ServeHTTP(w, r)
	})
}
