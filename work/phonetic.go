package work

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"words/helper"
	"words/model"

	"github.com/PuerkitoBio/goquery"
)

func AddPhonetic(p *model.Pool) {
	re, err := p.NoPhonetic(make(map[string]string))
	if err != nil {
		log.Fatal(err)
	}

	arr, ok := re.([]map[string]string)
	if !ok {
		log.Fatal("re.([]map[string]string) !ok")
	}
	for i, v := range arr {
		if v["parent_id"] == "0" {
			continue
		}
		title := strings.TrimSpace(v["title"])
		println(i, title)
		PhoneticFromWeb(p, v["id"], title)
	}
}

func PhoneticFromWeb(p *model.Pool, id, title string) {
	url := fmt.Sprintf("https://www.bing.com/dict/search?q=%s&qs=n&form=Z9LH5&sp=-1&pq=&sc=0-0&sk=&cvid=A5A51A9C2B71453EB239BA422322C361", strings.ToLower(title))

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	us := doc.Find(".hd_prUS.b_primtxt").Text()
	uk := doc.Find(".hd_pr.b_primtxt").Text()
	uss := strings.Split(us, "[")
	uks := strings.Split(uk, "[")
	if len(uss) > 1 {
		us = strings.Split(uss[1], "]")[0]
	} else {
		us = ""
	}
	if len(uks) > 1 {
		uk = strings.Split(uks[1], "]")[0]
	} else {
		uk = ""
	}
	println(us, uk)
	_, err = p.EditWord(map[string]string{"id": id, "us": helper.Translate(us), "uk": helper.Translate(uk)})
	if err != nil {
		log.Fatal(err)
	}
}
