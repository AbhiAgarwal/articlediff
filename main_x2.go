func (c *appContext) scrapeArticleHandler(w http.ResponseWriter, r *http.Request) {
    params := context.Get(r, "params").(httprouter.Params)
    var article articlesModel.ArticleResource
    url := params.ByName("url")
    if url != "" {
        doc, err := goquery.NewDocument(url)
        if err != nil {
            log.Fatal(err)
        }
        article.Data.Title = doc.Find("h1").Text()
        doc.Find("p").Each(func(i int, s *goquery.Selection) {
            article.Data.Article += (s.Text() + "\n")
        })
        fmt.Println(article.Data.Title)
        return
    } else {
        return
    }
}