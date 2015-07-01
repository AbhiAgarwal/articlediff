package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type Handler func(w http.ResponseWriter, r *http.Request)
type Resource map[string]Handler
type Actions struct {
	HandlerName string
	OneHandler  Handler
}

var resources = make(map[string]Resource)

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

func (r *router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

func NewRouter() *router {
	return &router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	resource, exists := resources[params.ByName("resource")]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action, exists := resource["list"]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action(w, r)
}

func OneHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	resource, exists := resources[params.ByName("resource")]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action, exists := resource["one"]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action(w, r)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	resource, exists := resources[params.ByName("resource")]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action, exists := resource["delete"]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action(w, r)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	resource, exists := resources[params.ByName("resource")]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action, exists := resource["put"]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action(w, r)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	resource, exists := resources[params.ByName("resource")]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action, exists := resource["post"]
	if exists == false {
		WriteError(w, ErrNotAcceptable)
		return
	}
	action(w, r)
}

func RegisterAction(ResourceName string, actions ...Actions) {
	for _, action := range actions {
		resourceName, actionName := ResourceName, action.HandlerName
		resource, exists := resources[resourceName]
		if exists == false {
			resource = make(Resource)
			resources[resourceName] = resource
		}
		resource[actionName] = action.OneHandler
	}
}
