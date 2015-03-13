package main

import (
	"log"
	"net/http"

	// "github.com/julienschmidt/httprouter"
	"github.com/mantishK/httprouter"
	"github.com/mantishK/wormhole/app/controller"
	"github.com/mantishK/wormhole/app/filters"
)

func main() {
	route()
}
func route() {
	//Create filters
	authenticateFilter := new(filters.Authenticate)
	authorizeFilter := new(filters.Authorize)

	//Create controller
	todoController := controller.Todo{}

	router := httprouter.New()

	//version
	ver := make([]string, 2, 5)
	ver[0] = ""
	ver[1] = "/v1"

	ns := "/api"

	for _, verName := range ver {
		router.GET(ns+verName+"/todo", todoController.Get, authenticateFilter, authorizeFilter)
		router.GET(ns+verName+"/todos", todoController.GetAllTodos, authenticateFilter, authorizeFilter)
		router.POST(ns+verName+"/todos", todoController.Add, authenticateFilter, authorizeFilter)
		router.DELETE(ns+verName+"/todos", todoController.Delete, authenticateFilter, authorizeFilter)
		router.PUT(ns+verName+"/todos", todoController.Update, authenticateFilter, authorizeFilter)
	}
	router.NotFound = http.FileServer(http.Dir("public")).ServeHTTP

	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
