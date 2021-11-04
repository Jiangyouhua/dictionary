package main

import "time"

type Catalogue struct {
	id         int
	parentId   int
	title      string
	info       string
	startDate  time.Time
	updateDate time.Time
}
