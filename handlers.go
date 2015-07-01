package main

import (
	"encoding/json"
	"net/http"
	// "log"
	// "fmt"

	articlesModel "github.com/abhiagarwal/articlediff/models"

	// "github.com/PuerkitoBio/goquery"
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
	updatearticle := Actions{
		HandlerName: "put",
		OneHandler:  c.updatearticleHandler,
	}
	RegisterAction("articles", listarticle, onearticle, updatearticle, deletearticle)

	// scrapearticle := Actions{
	// 	HandlerName: "post",
	// 	OneHandler:  c.scrapeArticleHandler,
	// }
	// RegisterAction("article", scrapearticle)
}

// func (c *appContext) scrapeArticleHandler(w http.ResponseWriter, r *http.Request) {
// 	params := context.Get(r, "params").(httprouter.Params)
// 	var article articlesModel.ArticleResource
// 	url := params.ByName("url")
// 	if url != "" {
// 		doc, err := goquery.NewDocument(url)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		article.Data.Title = doc.Find("h1").Text()
// 		doc.Find("p").Each(func(i int, s *goquery.Selection) {
// 			article.Data.Article += (s.Text() + "\n")
// 		})
// 		fmt.Println(article.Data.Title)
// 		return
// 	} else {
// 		return
// 	}
// }

func (c *appContext) articlesHandler(w http.ResponseWriter, r *http.Request) {
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	articles, err := repo.All()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(articles)
}

func (c *appContext) articleHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	article, err := repo.Find(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(article)
}

func (c *appContext) createarticleHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*articlesModel.ArticleResource)
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *appContext) deletearticleHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := articlesModel.ArticleRepo{c.db.C("articles")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}
