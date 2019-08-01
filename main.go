package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/mrlynn/shortener/config"
	"github.com/mrlynn/shortener/storage"
	"github.com/mrlynn/shortener/storage/mongodb"
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

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/go/:code", Redirect)
	router.GET("/info", Info)
	router.GET("/shortener", Get)
	router.POST("/shortener", Post)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router.Run(":" + port)

}

func Get(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}

func Post(c *gin.Context) {
	url := c.PostForm("url")

	if !govalidator.IsURL(url) {
		c.Writer.Write([]byte("URL is not valid"))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	genUrl, err := storage.SaveUrl(url)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
		"shortUrl": genUrl,
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
