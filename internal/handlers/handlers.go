package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dmitriy-zverev/blog-api/internal/db"
	"github.com/google/uuid"
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
		sendHttpMessage(w, "invalid params, check docs for more info", http.StatusBadRequest)
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

func (cfg *ApiConfig) PostsPutHandler(w http.ResponseWriter, req *http.Request) {
	postId, err := uuid.Parse(req.PathValue("postId"))
	if err != nil {
		sendHttpMessage(w, "invalid post id", http.StatusNotFound)
		return
	}

	oldPost, err := cfg.DB.GetPost(context.Background(), postId)
	if err != nil {
		sendHttpMessage(w, err.Error(), http.StatusNotFound)
		return
	}

	params, err := parseHttpJson[post](req.Body)
	if err != nil {
		sendHttpMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidPostParams(params) && !isValidPostTitleLength(params) && !isValidPostContentLength(params) {
		sendHttpMessage(w, "invalid params, check docs for more info", http.StatusBadRequest)
		return
	}

	arg := db.UpdatePostParams{
		Title:    params.Title,
		Content:  params.Content,
		Category: params.Category,
		Tags:     params.Tags,
		ID:       oldPost.ID,
	}

	newPost, err := cfg.DB.UpdatePost(context.Background(), arg)
	if err != nil {
		sendHttpMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendHttpMessage(w, newPost, http.StatusOK)
	log.Printf("Post (id: %s) updated at %v\n", newPost.ID, newPost.Updatedat)
}

func (cfg *ApiConfig) PostsDeleteHandler(w http.ResponseWriter, req *http.Request) {
	postId, err := uuid.Parse(req.PathValue("postId"))
	if err != nil {
		sendHttpMessage(w, "invalid post id", http.StatusNotFound)
		return
	}

	if err := cfg.DB.DeletePost(context.Background(), postId); err != nil {
		sendHttpMessage(w, "internal server error", http.StatusInternalServerError)
		return
	}

	sendHttpMessage(w, "", http.StatusNoContent)
	log.Printf("Post (id: %s) deleted\n", postId)
}
