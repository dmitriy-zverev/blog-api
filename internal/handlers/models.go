package handlers

import "github.com/dmitriy-zverev/blog-api/internal/db"

type ApiConfig struct {
	Build string
	Port  string
	DB    *db.Queries
}
