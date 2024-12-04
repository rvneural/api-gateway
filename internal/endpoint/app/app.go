package app

import (
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	router   *gin.Engine
	apiGroup *gin.RouterGroup
}

func New() *Endpoint {
	router := gin.Default()
	api := router.Group("/api")
	return &Endpoint{
		router:   router,
		apiGroup: api,
	}
}

func (e *Endpoint) SetApiKey(key string) {
	if key == "" {
		return
	}
	e.apiGroup.Use(func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != key {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	})
}

func (e *Endpoint) AddAPIEndpoint(pattern string, url string) {
	if pattern == "" || url == "" {
		return
	}
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/ + pattern"
	}
	e.apiGroup.Any(pattern, e.createReverseProxy(url))
}

func (e *Endpoint) AddStaticEndpoint(pattern string, path string) {
	if pattern == "" || path == "" {
		return
	}
	e.router.GET(pattern, func(c *gin.Context) {
		c.File(path)
	})
}

func (e *Endpoint) Start() {
	e.router.Run(":8000")
}

func (e *Endpoint) createReverseProxy(targetURL string) func(*gin.Context) {
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
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/media")

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
