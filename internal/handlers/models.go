package handlers

import "github.com/dmitriy-zverev/blog-api/internal/db"

type ApiConfig struct {
	Build string
	Port  string
	DB    *db.Queries
}

type post struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}
