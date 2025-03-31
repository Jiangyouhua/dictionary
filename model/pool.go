package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"words/helper"

	_ "github.com/go-sql-driver/mysql"
)

type Pool struct {
	db *sql.DB
}

type Router func(map[string]string) (interface{}, error)

func NewPool(isRoot bool) (*Pool, error) {
	s := "root:jiangyouhua@tcp(127.0.0.1:3306)/words"
	db, err := sql.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	return &Pool{db: db}, nil
}

func Routing(handle string, data map[string]string) []byte {
	p, err := NewPool(true)
	if err != nil {
		return RespondToBytes(1, err.Error(), nil)
	}
	defer p.db.Close()

	router := map[string]Router{
		"fetchBook":  p.FetchBook,
		"fetchWord":  p.FetchWord,
		"editBook":   p.EditBook,
		"editWord":   p.EditWord,
		"noPhonetic": p.NoPhonetic,
	}
	f, ok := router[handle]
	if !ok {
		return RespondToBytes(1, fmt.Sprintf("No router named '%s'", handle), nil)
	}

	re, err := f(data)
	if err != nil {
		return RespondToBytes(2, err.Error(), data)
	}
	return RespondToBytes(0, "", re)
}

func (p *Pool) FetchBook(data map[string]string) (interface{}, error) {
	return p.tree("book", data)
}

func (p *Pool) FetchWord(data map[string]string) (interface{}, error) {
	return p.line("word", data)
}

func (p *Pool) EditBook(data map[string]string) (interface{}, error) {
	return p.edit("book", data)
}

func (p *Pool) EditWord(data map[string]string) (interface{}, error) {
	return p.edit("word", data)
}

func (p *Pool) NoPhonetic(data map[string]string) (interface{}, error) {
	query := "select id, title from word where uk is null and title not like '% %' and title not like '% %'"
	return p.query(query)
}

func (p *Pool) UpdateBookLinkWord(title, key, value string) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE temp SET `%s`=`%s`+%s WHERE `title`='%s'", key, key, value, title)
	return p.db.Exec(query)
}

func (p *Pool) InsertOrUpdate(table string, data []map[string]string, insertColumns []string, updateColumns []string) (sql.Result, error) {
	var (
		values   = make([]string, 0, len(data))
		settings = make([]string, 0, len(updateColumns))
	)

	for _, col := range updateColumns {
		s := fmt.Sprintf("`%s`=VALUES(`%s`)", col, col)
		settings = append(settings, s)
	}

	for _, value := range data {
		arr := make([]string, len(insertColumns))
		for i, v := range insertColumns {
			arr[i] = helper.Translate(value[v])
		}
		row := fmt.Sprintf("('%s')", strings.Join(arr, "', '"))
		values = append(values, row)
	}

	query := fmt.Sprintf("INSERT INTO %s (`%s`) VALUES %s ON DUPLICATE KEY UPDATE %s", table, strings.Join(insertColumns, "`, `"), strings.Join(values, ","), strings.Join(settings, ","))

	return p.db.Exec(query)
}

func (p *Pool) line(table string, data map[string]string) (interface{}, error) {
	// 格式化查询条件。
	w := condition(data, false)
	query := fmt.Sprintf(`SELECT * FROM %s%s`, table, w)
	return p.query(query)
}

func (p *Pool) tree(table string, data map[string]string) (interface{}, error) {
	w := condition(data, true)
	query := fmt.Sprintf(`
	WITH RECURSIVE temp (id, parent_id, title, info, seat, state, path, level) AS (
  		SELECT id, parent_id, title, info, seat, state, CONCAT(100000000 + seat, title) as path, 0 as level FROM %s WHERE parent_id = 0
  		UNION ALL
  		SELECT c.id, c.parent_id, c.title, c.info, c.seat, c.state, CONCAT(cp.path, ' > ', 100000000 + c.seat, c.title) as path, cp.level+1 as level FROM temp AS cp JOIN %s AS c ON cp.id = c.parent_id
	)
	SELECT * FROM temp %s ORDER BY path;
	`, table, table, w)
	return p.query(query)
}

func (p *Pool) query(query string) ([]map[string]string, error) {
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	arr := make([]map[string]string, 0)
	for rows.Next() {
		var (
			temps  = make([]interface{}, len(columns))
			values = make([]sql.NullString, len(columns))
			row    = make(map[string]string)
		)
		for i := 0; i < len(columns); i++ {
			temps[i] = &(values[i])
		}
		if err = rows.Scan(temps...); err != nil {
			return nil, err
		}
		for k, v := range columns {
			row[v] = values[k].String
		}
		arr = append(arr, row)
	}
	return arr, nil
}

func (p *Pool) edit(table string, data map[string]string) (sql.Result, error) {
	var (
		columns  = make([]string, 0, len(data))
		values   = make([]string, 0, len(data))
		settings = make([]string, 0, len(data))
	)

	id := ""
	for k, v := range data {
		columns = append(columns, fmt.Sprintf("`%s`", k))
		values = append(values, fmt.Sprintf("'%s'", v))
		if k == "id" {
			id = v
			continue
		}
		settings = append(settings, fmt.Sprintf("`%s`='%s'", k, v))
	}
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %s`, table, strings.Join(settings, ", "), id)
	if id == "0" {
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, strings.Join(columns, ", "), strings.Join(values, ", "))
	}
	println(query)
	return p.db.Exec(query)
}

func condition(data map[string]string, isTree bool) string {
	wheres := make([]string, 0, len(data))
	pages := []int{1000, 0}
	order := ""

	for k, v := range data {
		switch k {
		case "pageNumber":
			if i, err := strconv.Atoi(v); err == nil {
				pages[0] = i
			}
		case "pageIndex":
			if i, err := strconv.Atoi(v); err == nil {
				pages[1] = i
			}
		case "orderColumn":
			order = fmt.Sprintf(" ORDER BY `%s` %s", k, v)
		default:
			wheres = append(wheres, fmt.Sprintf("`%s`='%s'", k, v))
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = fmt.Sprintf(" WHERE %s", strings.Join(wheres, " AND "))
	}
	if isTree {
		return where
	}
	limit := fmt.Sprintf(" LIMIT %v, %v", pages[0]*pages[1], pages[0])
	return where + order + limit
}
