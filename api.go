package main

import (
	"net/http"

	articlesModel "github.com/abhiagarwal/articlediff/models"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := appContext{session.DB("test")}
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler, acceptHandler)
	router := NewRouter()
	appC.register()

	router.Get("/:resource", commonHandlers.ThenFunc(ListHandler))
	router.Get("/:resource/:id", commonHandlers.ThenFunc(OneHandler))
	router.Delete("/:resource/:id", commonHandlers.ThenFunc(DeleteHandler))
	router.Put("/:resource/:id", commonHandlers.Append(contentTypeHandler, bodyHandler(articlesModel.ArticleResource{})).ThenFunc(UpdateHandler))
	router.Post("/:resource", commonHandlers.Append(contentTypeHandler, bodyHandler(articlesModel.ArticleResource{})).ThenFunc(CreateHandler))
	http.ListenAndServe(":8080", router)
}
