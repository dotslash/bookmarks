package main

import "net/http"

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Redirect",
		[]string{"GET"},
		"/red/{redId}",
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
}
