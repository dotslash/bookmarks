package main

import "net/http"

// Route contains information about the route and
// also the route's http handler/
type Route struct {
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

// Routes is an array of routes
// TODO(dotslash): Remove this.
type Routes []Route

var routes = Routes{
	Route{
		"Redirect",
		[]string{"GET"},
		"/red/{redId}",
		Redirect,
	},
	Route{
		"Redirect",
		[]string{"GET"},
		"/r/{redId}",
		Redirect,
	},
	Route{
		"ActionView",
		[]string{"GET", "POST"},
		"/actions/view",
		ActionView,
	},
	Route{
		"ActionAdd",
		[]string{"GET", "POST"},
		"/actions/add",
		ActionAdd,
	},
	Route{
		"ActionDel",
		[]string{"GET", "POST"},
		"/actions/delete",
		ActionDel,
	},
	Route{
		"ActionDel",
		[]string{"GET", "POST"},
		"/actions/delete",
		ActionDel,
	},
	Route{
		"ActionUpdate",
		[]string{"GET", "POST"},
		"/actions/update",
		ActionUpdate,
	},
	Route{
		"ActionLookup",
		[]string{"GET", "POST"},
		"/actions/lookup",
		ActionLookup,
	},
	Route{
		"ActionRevLookup",
		[]string{"GET", "POST"},
		"/actions/revlookup",
		ActionRevLookup,
	},
}
