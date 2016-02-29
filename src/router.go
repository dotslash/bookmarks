package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler
		log.Println(route)
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	return router
}
