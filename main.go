package main

import (
	"net/http"
	"strings"

	"log"

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

	repository := mongodb.NewMongoRepository(cfg.Mongo.URI, cfg.Mongo.DB, cfg.Mongo.Collection)

	storage.SetStorage(repository)

	router := gin.Default()

	router.LoadHTMLFiles("./static/index.html")

	router.GET("/redirect/:code", Redirect)
	router.GET("/shortener", Get)
	router.POST("/shortener", Post)

	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}

func Get(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Post(c *gin.Context) {
	url := c.PostForm("url")

	code, err := storage.SaveUrl(url)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"shortUrl": "http://localhost:8080/redirect/" + code,
	})
}

func Redirect(c *gin.Context) {
	code := c.Param("code")

	log.Println(code)

	url, err := storage.GetURL(code)

	if err != nil {
		log.Println("ERROR", err)
		c.Writer.WriteHeader(http.StatusNotFound)
		c.Writer.Write([]byte("URL Not Found"))
		return
	}

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	log.Println(url)

	http.Redirect(c.Writer, c.Request, url, http.StatusMovedPermanently)
}
