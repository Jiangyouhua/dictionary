package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"words/model"
	"words/work"
)

func main() {
	// service()
	format()
}

func format() {
	p, err := model.NewPool(true)
	if err != nil {
		log.Fatal(err)
	}
	// work.BookClassification(p, "book")
	// work.WordSupplement(p, 2)
	// work.DirSupplement(p, false)
	// work.SyllableToDatabase(p)
	// work.IntoDatabasesWith70w(p)
	work.AddPhonetic(p)
	// work.PhoneticFromWeb(p, "1", "anguish")
}

func service() {
	http.HandleFunc("/router", handle)
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", body)

	var m map[string]string
	json.Unmarshal([]byte(body), &m)

	h := m["handle"]
	delete(m, "handle")
	if len(h) == 0 {
		http.NotFound(w, r)
		return
	}
	b := model.Routing(h, m)
	w.Write(b)
}
