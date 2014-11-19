package router

import (
	"net/http"
	"strings"

	"github.com/mantishK/wormhole/app/filters"
	"github.com/mantishK/wormhole/app/views"
)

type router struct {
	routeMethodMap  map[string]map[string]func(http.ResponseWriter, *http.Request, map[string]interface{})
	routerFilterMap map[string]map[string][]filters.Filterable
	nameSpace       string
	publicNameSpace string
}

func New() *router {
	newRouter := new(router)
	newRouter.routeMethodMap = make(map[string]map[string]func(http.ResponseWriter, *http.Request, map[string]interface{}))
	newRouter.routerFilterMap = make(map[string]map[string][]filters.Filterable)
	newRouter.nameSpace = "/api"
	newRouter.publicNameSpace = ""
	return newRouter
}

func NewWithNameSpace(nameSpace, publicNameSpace string) *router {
	newRouter := New()
	if strings.TrimSpace(nameSpace) == "" {
		newRouter.nameSpace = ""
	} else {
		newRouter.nameSpace = "/" + nameSpace
	}
	if strings.TrimSpace(publicNameSpace) == "" {
		newRouter.publicNameSpace = ""
	} else {
		newRouter.publicNameSpace = "/" + publicNameSpace
	}
	return newRouter
}

func (r *router) setMethodMap(httpMethod, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice []filters.Filterable) {
	if !strings.HasSuffix(route, "/") {
		route = route + "/"
	}
	if r.routeMethodMap[route] == nil {
		r.routeMethodMap[route] = make(map[string]func(http.ResponseWriter, *http.Request, map[string]interface{}))
	}
	if len(filterSlice) != 0 {
		if r.routerFilterMap[route] == nil {
			r.routerFilterMap[route] = make(map[string][]filters.Filterable)
			if r.routerFilterMap[route][httpMethod] == nil {
				r.routerFilterMap[route][httpMethod] = make([]filters.Filterable, 0, 5)
			}
		}
		for i := range filterSlice {
			r.routerFilterMap[route][httpMethod] = append(r.routerFilterMap[route][httpMethod], filterSlice[i])
		}
	}
	r.routeMethodMap[route][httpMethod] = newControllerMethod
}

func (r *router) Get(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("GET", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) Post(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("POST", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) Put(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("PUT", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) Delete(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("DELETE", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) Head(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("HEAD", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) Options(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("OPTIONS", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) Trace(versionName, route string, newControllerMethod func(http.ResponseWriter, *http.Request, map[string]interface{}),
	filterSlice ...filters.Filterable) {
	r.setMethodMap("TRACE", versionName+r.nameSpace+route, newControllerMethod, filterSlice)
}

func (r *router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	method := req.Method
	requestURI := strings.Split(req.RequestURI, "?")[0]
	if !strings.HasSuffix(requestURI, "/") {
		requestURI = requestURI + "/"
	}
	matchesPublic := strings.HasPrefix(requestURI, r.publicNameSpace)
	matchesApi := strings.HasPrefix(requestURI, r.nameSpace)
	if (matchesPublic && !matchesApi) || (matchesPublic && len(r.publicNameSpace) > len(r.nameSpace)) {
		http.ServeFile(writer, req, "public/"+strings.TrimLeft(req.URL.Path, r.publicNameSpace))
		return
	}
	filterData := make(map[string]interface{})
	filterReturnVal := r.executeFilters(method, requestURI, writer, req, filterData)
	if filterReturnVal == true {
		if r.routeMethodMap[requestURI][method] != nil {
			r.routeMethodMap[requestURI][method](writer, req, filterData)
		} else {
			view := views.NewView(writer)
			view.RenderHttpError("404!! Resource Not Found", 404)
		}
	}
}

func (r *router) executeFilters(method, requestURI string, writer http.ResponseWriter, req *http.Request, filterData map[string]interface{}) bool {
	filtersSlice := r.routerFilterMap[requestURI][method]
	for i := range filtersSlice {
		returnVal := r.routerFilterMap[requestURI][method][i].Filter(writer, req, filterData)
		if returnVal == false {
			return false
		}
	}
	return true
}
