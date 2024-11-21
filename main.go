package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	Text2TextURL := os.Getenv("Text2TextURL")
	Audio2TextURL := os.Getenv("Audio2TextURL")
	Text2ImageURL := os.Getenv("Text2ImageURL")
	RSSUrl := os.Getenv("RSSU")
	DBURL := os.Getenv("DBURL")
	API_KEY := os.Getenv("API_KEY")

	textProxy := createReverseProxy(Text2TextURL)
	audioProxy := createReverseProxy(Audio2TextURL)
	imageProxy := createReverseProxy(Text2ImageURL)
	rssProxy := createReverseProxy(RSSUrl)
	dbProxy := createReverseProxy(DBURL)

	r := gin.Default()
	group := r.Group("/api")
	if API_KEY != "" {
		group.Use(authenticate(API_KEY))
	}
	group.Any("/text/*path", textProxy)
	group.Any("/audio/*path", audioProxy)
	group.Any("/image/*path", imageProxy)
	group.Any("/db/*path", dbProxy)
	group.Any("/rss/*path", rssProxy)

	r.StaticFile("/docs.html", "./docs.html")

	r.GET("/docs", func(c *gin.Context) {
		c.File("docs.html")
	})

	log.Fatal(r.Run(":8000"))
}

func authenticate(APIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != APIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func createReverseProxy(targetURL string) func(*gin.Context) {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		log.Printf("Request received for %s", c.Request.URL.Path)

		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/text")
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/audio")
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/image")
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/db")
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/rss")

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
