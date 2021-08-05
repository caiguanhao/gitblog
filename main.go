package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caiguanhao/gitdb"
	"github.com/gin-gonic/gin"
	"github.com/gopsql/goconf"
)

var (
	configFileLocation string
)

func getConfigs() *Configs {
	c := Configs{}
	content, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return &c
	}
	goconf.Unmarshal(content, &c)
	return &c
}

func main() {
	defaultConfigFile := ".gitblog.go"
	if home, _ := os.UserHomeDir(); home != "" {
		defaultConfigFile = filepath.Join(home, defaultConfigFile)
	}
	flag.StringVar(&configFileLocation, "c", defaultConfigFile, "location of the config file")
	toUpdateConfigs := flag.Bool("C", false, "create or update config file")
	flag.Parse()

	configs := getConfigs()
	if *toUpdateConfigs {
		updateConfigs(configs)
		return
	}

	db := gitdb.NewDB(configs.Remote, configs.Local)
	db.SetSSHKey("git", configs.SSHPrivateKey, configs.SSHPrivateKeyPassword)
	db.SetUser(configs.UserName, configs.UserEmail)
	db.MustInit()

	modelPost := db.NewCollection("posts.js")
	modelPost.JSONPCallbackName = "__renderPosts"

	api := API{
		db:        db,
		modelPost: modelPost,
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(api.handleError)
	g := r.Group("/api")
	g.GET("/status", api.getStatus)
	g.POST("/push", api.push)
	g.GET("/posts", api.getPosts)
	g.GET("/posts/:id", api.showPost)
	g.POST("/posts", api.createPost)
	g.PUT("/posts/:id", api.updatePost)
	g.DELETE("/posts/:id", api.destroyPost)

	address := configs.Address
	if address == "" {
		address = "127.0.0.1:8080"
	}
	r.Run(address)
}
