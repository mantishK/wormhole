package main

import (
	"log"
	"net/http"

	"github.com/mantishK/wormhole/app/controller"
	"github.com/mantishK/wormhole/app/filters"
	"github.com/mantishK/wormhole/router"
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

	//version
	version := make([]string, 5, 5)
	version[0] = ""
	version[1] = "/v1"

	//Add new versions as and when required
	//version[2] = "/v2"

	//Create Router
	myRouter := router.New()

	//Create router with name spaces for api and public files
	//myRouter := router.NewNewWithNameSpace("api","app")

	//route
	for _, versionName := range version {
		//todos
		myRouter.Get(versionName, "/todo", todoController.Get, authenticateFilter)
		myRouter.Get(versionName, "/todos", todoController.GetAllTodos, authenticateFilter)
		myRouter.Post(versionName, "/todos", todoController.Add, authenticateFilter)
		myRouter.Delete(versionName, "/todos", todoController.Delete, authenticateFilter)
		myRouter.Put(versionName, "/todos", todoController.Update, authenticateFilter, authorizeFilter)

	}

	//New versions to be easily overridden
	// myRouter.Get(version[1], "/todo", noteController.GetV1, authenticateFilter)

	http.Handle("/", myRouter)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
