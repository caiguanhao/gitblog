package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caiguanhao/gitdb"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type (
	API struct {
		db        *gitdb.DB
		modelPost *gitdb.Collection
	}
)

func (api API) NewObject(post Post) *gitdb.Object {
	modelBody := api.db.NewObject("posts/" + post.Id + ".js")
	modelBody.JSONPCallbackName = "__renderPost"
	return modelBody
}

func (api API) getPosts(c *gin.Context) {
	var posts []Post
	api.modelPost.MustRead(&posts)
	c.JSON(200, posts)
}

func (api API) showPost(c *gin.Context) {
	var posts []Post
	api.modelPost.MustRead(&posts)
	for _, post := range posts {
		if post.Id == c.Param("id") {
			var full PostFull
			modelBody := api.NewObject(post)
			modelBody.MustRead(&full)
			c.JSON(200, full)
			return
		}
	}
	c.Status(404)
}

func (api API) createPost(c *gin.Context) {
	var postReq PostRequest
	if err := c.ShouldBindJSON(&postReq); err != nil {
		panic(err)
		return
	}
	var posts []Post
	api.modelPost.MustRead(&posts)
	now := time.Now()
	n := now.Round(time.Second)
	post := Post{
		Id:        strconv.FormatInt(now.UnixNano(), 10),
		Title:     postReq.Title,
		CreatedAt: &n,
		UpdatedAt: &n,
	}
	posts = append([]Post{post}, posts...)
	api.modelPost.MustWrite(posts)
	modelBody := api.NewObject(post)
	full := PostFull{
		post,
		postReq.Body,
	}
	modelBody.MustWrite(full)
	api.db.MustAdd(api.modelPost.Path, modelBody.Path)
	api.db.MustCommit()
	c.JSON(200, full)
}

func (api API) updatePost(c *gin.Context) {
	var postReq PostRequest
	if err := c.ShouldBindJSON(&postReq); err != nil {
		panic(err)
		return
	}
	var posts []Post
	var full PostFull
	api.modelPost.MustRead(&posts)
	api.modelPost.MustWrite(posts, func(post *Post) *Post {
		if post.Id == c.Param("id") {
			post.Title = postReq.Title
			now := time.Now().Round(time.Second)
			post.UpdatedAt = &now

			modelBody := api.NewObject(*post)
			full = PostFull{
				*post,
				postReq.Body,
			}
			modelBody.MustWrite(full)
			api.db.MustAdd(modelBody.Path)
		}
		return post
	})
	api.db.MustAdd(api.modelPost.Path)
	api.db.MustCommit()
	c.JSON(200, full)
}

func (api API) destroyPost(c *gin.Context) {
	var posts []Post
	api.modelPost.MustRead(&posts)
	api.modelPost.MustWrite(posts, func(post *Post) *Post {
		if post.Id == c.Param("id") {
			modelBody := api.NewObject(*post)
			modelBody.MustDelete()
			api.db.MustAdd(modelBody.Path)
			return nil // delete
		}
		return post
	})
	api.db.MustAdd(api.modelPost.Path)
	api.db.MustCommit(fmt.Sprint("destroyed post with id=", c.Param("id")))
	c.Status(204)
}

func (api API) getStatus(c *gin.Context) {
	commits := api.db.MustUnpushedCommits()
	c.JSON(200, gin.H{"UnpushedCommits": len(commits)})
}

func (api API) push(c *gin.Context) {
	api.db.MustPush()
	c.Status(204)
}

func (api API) handleError(c *gin.Context) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if err, ok := r.(validator.ValidationErrors); ok {
			var errs []gin.H
			for _, e := range err {
				errs = append(errs, gin.H{
					"FullName": e.Namespace(),
					"Name":     e.Field(),
					"Kind":     e.Kind().String(),
					"Type":     e.Tag(),
					"Param":    e.Param(),
				})
			}
			c.JSON(400, gin.H{"Errors": errs})
			return
		}
		if err, ok := r.(error); ok {
			c.JSON(400, gin.H{"Error": err.Error()})
			return
		}
		c.AbortWithStatus(500)
	}()
	c.Next()
}
