package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Readability struct {
	Author        string      `json:"author"`
	Content       string      `json:"content"`
	DatePublished string      `json:"date_published"`
	Dek           interface{} `json:"dek"`
	Direction     string      `json:"direction"`
	Domain        string      `json:"domain"`
	Excerpt       string      `json:"excerpt"`
	LeadImageURL  interface{} `json:"lead_image_url"`
	NextPageID    interface{} `json:"next_page_id"`
	RenderedPages int         `json:"rendered_pages"`
	ShortURL      string      `json:"short_url"`
	Title         string      `json:"title"`
	TotalPages    int         `json:"total_pages"`
	URL           string      `json:"url"`
	WordCount     int         `json:"word_count"`
}

var readabilityToken = "de063260cde5ec0147a27a6999ab720b392ed952"

func getReadabilityData(url string) Readability {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://readability.com/api/content/v1/parser?url="+url+"/&token="+readabilityToken, nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	b := bytes.NewBufferString(string(body))
	decoder := json.NewDecoder(b)
	var result Readability
	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalln(err)
	}
	return result
}

func main() {
	data := getReadabilityData("http://www.huffingtonpost.com/2015/06/30/greece-imf-debt_n_7658084.html")
	fmt.Println(data.Content)
}
