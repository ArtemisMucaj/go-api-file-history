package main

import "log"
import "net/http"

func main() {
	ctx, err := NewAppContext()
	if err != nil {
		log.Fatal("Could not connect to mongodb -> ", err)
	}
	defer ctx.Close()
	routes := Routes{
		Route{"index", "GET", "/", ctx.Index, false},
		Route{"dev", "GET", "/dev", ctx.Dev, false},
		Route{"login", "POST", "/login", ctx.Login, false},
		Route{"signin", "POST", "/signin", ctx.Signin, false},
		Route{"getfile", "POST", "/getfile", ctx.GetFile, true},
		Route{"createfile", "POST", "/createfile", ctx.CreateFile, true},
		Route{"savefile", "POST", "/savefile", ctx.SaveFile, true}}
	router := NewRouter(ctx, routes)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/")))
	log.Fatal(http.ListenAndServe(":8080", router))
}
