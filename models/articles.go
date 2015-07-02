package models

import (
  "time"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type article struct {
  Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Title    string        `json:"title"`
  Article  string        `json:"article"`
  URL      string        `json:"url"`
  Date     time.Time     `json:"date"`
}

type ArticlesCollection struct {
  Data []article `json:"data"`
}

type ArticleResource struct {
  Data article `json:"data"`
}

type ArticleRepo struct {
  Coll *mgo.Collection
}

func (r *ArticleRepo) All() (ArticlesCollection, error) {
  result := ArticlesCollection{[]article{}}
  err := r.Coll.Find(nil).All(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *ArticleRepo) Find(id string) (ArticleResource, error) {
  result := ArticleResource{}
  err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *ArticleRepo) Create(article *article) error {
  id := bson.NewObjectId()
  _, err := r.Coll.UpsertId(id, article)
  if err != nil {
    return err
  }

  article.Id = id

  return nil
}

func (r *ArticleRepo) Update(article *article) error {
  err := r.Coll.UpdateId(article.Id, article)
  if err != nil {
    return err
  }

  return nil
}

func (r *ArticleRepo) Delete(id string) error {
  err := r.Coll.RemoveId(bson.ObjectIdHex(id))
  if err != nil {
    return err
  }

  return nil
}