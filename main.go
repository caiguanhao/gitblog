package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/caiguanhao/gitdb"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
)

var (
	configFileLocation string
	noConfigsApi       bool

	serverQuitChan chan os.Signal

	frontendFS http.FileSystem
)

func main() {
	defaultConfigFile := ".gitblog.go"
	if home, _ := os.UserHomeDir(); home != "" {
		defaultConfigFile = filepath.Join(home, defaultConfigFile)
	}
	flag.StringVar(&configFileLocation, "c", defaultConfigFile, "location of the config file")
	toUpdateConfigs := flag.Bool("C", false, "create or update config file")
	flag.BoolVar(&noConfigsApi, "no-configs", false, "disable the configs api")
	flag.Parse()

	if *toUpdateConfigs {
		updateConfigs(getConfigs())
		return
	}

	for {
		startServer()
	}
}

func restartServer() {
	serverQuitChan <- syscall.Signal(0xa)
}

func startServer() {
	configs := getConfigs()

	db := gitdb.NewDB(configs.Remote, configs.Local)
	db.SetSSHKey("git", configs.SSHPrivateKey, configs.SSHPrivateKeyPassword)
	db.SetUser(configs.UserName, configs.UserEmail)
	if configs.Remote != "" && configs.Local != "" {
		db.MustInit()
	}

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
	if noConfigsApi == false {
		g.GET("/configs", api.getConfigs)
		g.POST("/configs", api.updateConfigs)
	}
	g.GET("/status", api.getStatus)
	g.POST("/push", api.push)
	g.GET("/posts", api.getPosts)
	g.GET("/posts/:id", api.showPost)
	g.POST("/posts", api.createPost)
	g.PUT("/posts/:id", api.updatePost)
	g.DELETE("/posts/:id", api.destroyPost)

	if frontendFS != nil {
		r.NoRoute(func(c *gin.Context) {
			testResp := httptest.NewRecorder()
			test, _ := gin.CreateTestContext(testResp)
			test.Request = c.Request
			test.FileFromFS(test.Request.URL.Path, frontendFS)
			if testResp.Code == 404 {
				c.Request.URL.Path = "/"
			}
			c.FileFromFS(c.Request.URL.Path, frontendFS)
		})
	}

	address := configs.GetAddress()
	var portNumber int
	if tcpAddr, err := net.ResolveTCPAddr("tcp", address); err == nil {
		portNumber = tcpAddr.Port
	}

	// open configs page for the first time
	if noConfigsApi == false && !hasConfigs() && frontendFS != nil && portNumber > 0 {
		browser.OpenURL(fmt.Sprintf("http://127.0.0.1:%d/configs", portNumber))
	}

	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}
	go func() {
		if portNumber > 0 {
			log.Printf("Listening %s (http://127.0.0.1:%d)", address, portNumber)
		}
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("Listen error:", err)
		}
	}()
	serverQuitChan = make(chan os.Signal)
	signal.Notify(serverQuitChan, syscall.SIGINT, syscall.SIGTERM)
	signal := <-serverQuitChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("Error shutting down server:", err)
	}
	log.Println("Server shut down successfully")
	if signal == syscall.Signal(0xa) {
		return
	}
	os.Exit(0)
}
