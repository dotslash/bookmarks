package main

import (
	"net/http"

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

func getRoutes() []routeStruct {
	getOnly := []string{"GET"}
	getAndPost := []string{"GET", "POST"}
	return []routeStruct{
		routeStruct{"Redirect", getOnly, "/red/{redId}", Redirect},
		routeStruct{"Redirect", getOnly, "/r/{redId}", Redirect},
		routeStruct{"ActionView", getAndPost, "/actions/view", ActionView},
		routeStruct{"ActionAdd", getAndPost, "/actions/add", ActionAdd},
		routeStruct{"ActionDel", getAndPost, "/actions/delete", ActionDel},
		routeStruct{"ActionDel", getAndPost, "/actions/delete", ActionDel},
		routeStruct{"ActionUpdate", getAndPost, "/actions/update", ActionUpdate},
		routeStruct{"ActionLookup", getAndPost, "/actions/lookup", ActionLookup},
		routeStruct{
			"ActionRevLookup", getAndPost, "/actions/revlookup", ActionRevLookup},
	}
}

// NewRouter returns a new mux.Router that handles all incoming http requests.
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range getRoutes() {
		var handler http.Handler
		Log.Println(route)
		handler = route.HandlerFunc
		handler = HTTPLogger(handler, route.Name)
		router.Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	return router
}
