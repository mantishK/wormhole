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
	userController := controller.User{}

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
		myRouter.Get(versionName, "/todo", todoController.Get, authenticateFilter, authorizeFilter)
		myRouter.Get(versionName, "/todos", todoController.GetAllUserTodos, authenticateFilter, authorizeFilter)
		myRouter.Post(versionName, "/todos", todoController.Add, authenticateFilter, authorizeFilter)
		myRouter.Delete(versionName, "/todos", todoController.Delete, authenticateFilter, authorizeFilter)
		myRouter.Put(versionName, "/todos", todoController.Update, authenticateFilter, authorizeFilter)

		//User
		myRouter.Post(versionName, "/user", userController.SaveUser, authenticateFilter)
		myRouter.Put(versionName, "/user", userController.UpdatePassword, authenticateFilter, authorizeFilter)
		myRouter.Get(versionName, "/user", userController.GetUser, authorizeFilter, authenticateFilter)

		//Check if user exists, could also be a Head request for user rather than Get
		myRouter.Get(versionName, "/user_name_exists", userController.UserNameExists, authenticateFilter)

		//Sessions
		myRouter.Post(versionName, "/session", userController.SignIn, authenticateFilter)
		myRouter.Delete(versionName, "/session", userController.SignOut, authenticateFilter, authorizeFilter)
		myRouter.Put(versionName, "/session", userController.RenewToken, authenticateFilter, authorizeFilter)

	}

	//New versions to be easily overridden
	// myRouter.Get(version[1], "/todo", noteController.GetV1, authenticateFilter)

	http.Handle("/", myRouter)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
