package main

import (
	"encoding/json"
	"time"
)

type (
	Post struct {
		Id        string
		Title     string
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}

	PostFull struct {
		Post
		Body string
	}

	PostRequest struct {
		Title string `binding:"required,lt=100"`
		Body  string `binding:"required,lt=10000"`
	}
)

func (post PostFull) GITDBMarshalJSON() []byte {
	j, _ := json.MarshalIndent(post, "", "  ")
	return j
}
