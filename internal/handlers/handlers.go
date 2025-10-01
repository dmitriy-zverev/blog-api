package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dmitriy-zverev/blog-api/internal/db"
)

func (cfg *ApiConfig) BaseHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) PostsPostHandler(w http.ResponseWriter, req *http.Request) {
	params, err := parseHttpJson[post](req.Body)
	if err != nil {
		sendHttpMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidPostParams(params) && !isValidPostTitleLength(params) && !isValidPostContentLength(params) {
		sendHttpMessage(w, "invalid params, check docs for more info", http.StatusInternalServerError)
		return
	}

	arg := db.CreatePostParams{
		Title:    params.Title,
		Content:  params.Content,
		Category: params.Category,
		Tags:     params.Tags,
	}

	post, err := cfg.DB.CreatePost(context.Background(), arg)
	if err != nil {
		sendHttpMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendHttpMessage(w, post, http.StatusCreated)
	log.Printf("Post (id: %s) created at %v\n", post.ID, post.Createdat)
}
