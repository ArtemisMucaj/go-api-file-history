package main

import "net/http"
import "github.com/gorilla/mux"

func NewRouter(ctx *AppContext, routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.Authenticate {
			handler = ctx.Authenticate(handler)
		}
		handler = ctx.SetDatabase(handler)
		handler = ctx.Logger(handler, route.Name)
		router.Methods(route.Method).Path(route.Pattern).
			Name(route.Name).Handler(handler)
	}
	return router
}
