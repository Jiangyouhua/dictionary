package work

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"words/helper"
	"words/model"
)

// 根据book文件夹向数据库写入书本分类。
func BookClassification(p *model.Pool, filePath string, id int) {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		println("ioutil.ReadDir", filePath)
		return
	}
	for _, v := range files {
		name := v.Name()
		if name[0] == '.' {
			continue
		}
		i := writeToDatabase(p, id, name)
		if v.IsDir() {
			BookClassification(p, fmt.Sprintf("%s/%s", filePath, name), i)
		}
	}
}

func writeToDatabase(p *model.Pool, parentId int, title string) int {
	println(parentId, title)
	m := make(map[string]string)
	m["id"] = "0"
	m["parent_id"] = strconv.Itoa(parentId)
	m["title"] = title
	re, err := p.EditBook(m)
	if err != nil {
		log.Fatal(err, re)
	}
	r, ok := re.(sql.Result)
	if !ok {
		log.Fatal("re.(sql.Result) !ok")
	}
	id, err := r.LastInsertId()
	if err != nil {
		log.Fatal(re, id)
	}
	return int(id)
}

// 从数据库获取书本，按照从简单到复杂，写入单词与释义。
func WordSupplement(p *model.Pool, method int) {
	filePath := "csv"
	if method > 0 {
		filePath = "txt"
	}
	re, err := p.FetchBook(nil)
	if err != nil {
		log.Fatal(err, re)
	}
	arr, ok := re.([]map[string]string)
	if !ok {
		log.Fatal("re.([]map[string]string) !ok")
	}
	for i, v := range arr {
		if v["parent_id"] == "0" {
			continue
		}
		file := strings.TrimSpace(v["title"])
		f := fmt.Sprintf(`%s/%s`, filePath, strings.Split(file, ".")[0])
		print(i, ",")
		switch method {
		case 0:
			infoToDatabase(p, f)
		case 1:
			titleToDatabase(p, f)
		default:
			wordLinkBook(p, v["id"], f)
		}
	}
}

func DirSupplement(p *model.Pool, isTitle bool) {
	filePath := "csv"
	if isTitle {
		filePath = "txt"
	}
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		println("ioutil.ReadDir", filePath)
		return
	}
	for i, v := range files {
		name := v.Name()
		if name[0] == '.' || v.IsDir() {
			continue
		}
		f := fmt.Sprintf(`%s/%s`, filePath, strings.Split(name, ".")[0])
		print(i, ",")
		if isTitle {
			titleToDatabase(p, f)
		} else {
			infoToDatabase(p, f)
		}
	}
}

func titleToDatabase(p *model.Pool, fileName string) {
	file, err := os.Open(fileName + ".txt")
	if err != nil {
		return
	}
	defer file.Close()

	data := make([]map[string]string, 0)
	r := bufio.NewScanner(file)
	for r.Scan() {
		line := helper.Translate(r.Text())
		if len(line) == 0 || strings.Index(line, "#") > -1 {
			continue
		}
		data = append(data, map[string]string{"title": line, "info": ""})
	}
	_, err = p.InsertOrUpdate("temp", data, []string{"title", "info"}, []string{"info"})
	if err != nil {
		println("InsertOrUpdate", file)
	}
	os.Remove(fileName + ".txt")
	// log.Println(re)
}

func infoToDatabase(p *model.Pool, fileName string) {
	file, err := os.Open(fileName + ".csv")
	if err != nil {
		return
	}
	defer file.Close()

	data := make([]map[string]string, 0)
	r := csv.NewReader(file)
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		title := helper.Translate(row[0])
		info := helper.Translate(row[1])
		if info == "无" {
			continue
		}
		// println("*", title, ":", info)
		data = append(data, map[string]string{"title": title, "info": info})
	}
	_, err = p.InsertOrUpdate("temp", data, []string{"title", "info"}, []string{"info"})
	if err != nil {
		println("InsertOrUpdate", file)
	}
	os.Remove(fileName + ".csv")
	// log.Println(re)
}

// 从word表，写入音节、音标。
func SyllableToDatabase(p *model.Pool) {
	for i := 0; i < 100; i++ {
		println(i)
		wordToTemp(p, strconv.Itoa(i))
	}
}

func wordToTemp(p *model.Pool, index string) {
	data := map[string]string{
		"pageIndex":  index,
		"pageNumber": "1000",
	}
	result, err := p.FetchWord(data)
	if err != nil {
		log.Println(err)
	}
	r, ok := result.([]map[string]string)
	if !ok || len(r) == 0 {
		log.Fatal(result, "result.([]map[string]string)")
	}
	re, err := p.InsertOrUpdate("temp", r, []string{"title", "info", "syllable", "uk", "us"}, []string{"info", "syllable", "uk", "us"})
	if err != nil {
		log.Fatal(err, re)
	}
	log.Println(re)
}

func wordLinkBook(p *model.Pool, bookId, fileName string) {
	file, err := os.Open(fileName + ".txt")
	if err != nil {
		return
	}
	defer file.Close()

	id, err := strconv.Atoi(bookId)
	if err != nil {
		return
	}

	divide := 62 // 2的63次方超出Mysql BIGINT UNSIGNED的最大值。
	column := fmt.Sprintf(`book%v`, id/divide)
	var value uint64 = 1 << (id % divide)
	v := strconv.FormatUint(value, 10)
	r := bufio.NewScanner(file)
	for r.Scan() {
		line := helper.Translate(r.Text())
		if len(line) == 0 || strings.Index(line, "#") > -1 {
			continue
		}
		if _, err := p.UpdateBookLinkWord(line, column, v); err != nil {
			log.Fatal(err)
		}
	}
	os.Remove(fileName + ".txt")
	// log.Println(re)
}
