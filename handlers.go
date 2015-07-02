package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	articlesModel "github.com/abhiagarwal/articlediff/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type appContext struct {
	db *mgo.Database
}

func (c *appContext) register() {
	listarticle := Actions{
		HandlerName: "list",
		OneHandler:  c.articlesHandler,
	}
	onearticle := Actions{
		HandlerName: "one",
		OneHandler:  c.articleHandler,
	}
	deletearticle := Actions{
		HandlerName: "delete",
		OneHandler:  c.deletearticleHandler,
	}
	createarticle := Actions{
		HandlerName: "post",
		OneHandler:  c.createarticleHandler,
	}
	updatearticle := Actions{
		HandlerName: "put",
		OneHandler:  c.updatearticleHandler,
	}
	RegisterAction("articles", listarticle, onearticle, createarticle, updatearticle, deletearticle)
}

func (c *appContext) articlesHandler(w http.ResponseWriter, r *http.Request) {
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	articles, err := repo.All()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(articles)
}

func (c *appContext) articleHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	article, err := repo.Find(params.ByName("id"))
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(article)
}

func (c *appContext) createarticleHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*articlesModel.ArticleResource)
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	doc, err := goquery.NewDocument(body.Data.URL)
	if err != nil {
		log.Fatal(err)
	}

	body.Data.Title = doc.Find("h1").Text()
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		body.Data.Article += html
	})
	body.Data.Date = time.Now()

	err = repo.Create(&body.Data)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (c *appContext) updatearticleHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*articlesModel.ArticleResource)
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	err := repo.Update(&body.Data)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *appContext) deletearticleHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}
