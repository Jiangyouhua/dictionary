package main

import (
	"database/sql"
	"log"
)

type Pool struct {
	db *sql.DB
}

func NewPool(isRoot bool) *Pool {
	s := "root:password@tcp(127.0.0.1:3306)/grammar"
	if !isRoot {
		s = "root:password@tcp(118.31.78.122:3306)/grammar"
	}
	db, err := sql.Open("mysql", s)
	if err != nil {
		log.Fatal(err)
	}
	return &Pool{db: db}
}

func (p *Pool) selectAllByTreeTable(parentId int) {
	sql := `
	WITH RECURSIVE category_path (id, title, path) AS
                   (
                       SELECT id, title, title as path
                       FROM category
                       WHERE parent_id = null
                       UNION ALL
                       SELECT c.id, c.title, CONCAT(cp.path, ' > ', c.title)
                       FROM category_path AS cp JOIN category AS c
                                                     ON cp.id = c.parent_id
                   )
SELECT * FROM category_path
ORDER BY path;
	`
	p.db.Query(sql)
}
