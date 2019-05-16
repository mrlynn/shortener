package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Dimitriy14/shortener/config"
	"github.com/Dimitriy14/shortener/storage"
	"github.com/Dimitriy14/shortener/storage/mongodb"
)

func main() {
	cfg, err := config.GetConfigFromJSON("config.json")

	if err != nil {
		log.Fatal(err)
	}

	client, err := mongodb.NewMongoClient(cfg.Mongo.URI)

	if err != nil {
		log.Fatal(err)
	}

	repository := mongodb.NewMongoRepository(cfg.Mongo.DB, cfg.Mongo.Collection, client)

	storage.SetStorage(repository)

	router := gin.Default()

	router.LoadHTMLFiles("./static/index.html")

	router.GET("/go/:code", Redirect)
	router.GET("/info", Info)
	router.GET("/shortener", Get)
	router.POST("/shortener", Post)

	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}

func Get(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Post(c *gin.Context) {
	url := c.PostForm("url")

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	genUrl, err := storage.SaveUrl(url)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"shortUrl": genUrl,
	})
}

func Redirect(c *gin.Context) {
	code := c.Param("code")

	url, err := storage.GetURL(code)

	if err != nil {
		log.Println("ERROR", err)
		c.Writer.WriteHeader(http.StatusNotFound)
		c.Writer.Write([]byte("URL Not Found"))
		return
	}

	http.Redirect(c.Writer, c.Request, url, http.StatusMovedPermanently)
}

func Info(c *gin.Context) {
	info, err := storage.GetInfo()

	if err != nil {
		log.Println("ERROR", err)
		c.Writer.WriteHeader(http.StatusNotFound)
		c.Writer.Write([]byte("URL Not Found"))
		return
	}

	c.JSON(http.StatusOK, info)
}
