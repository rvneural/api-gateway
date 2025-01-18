package app

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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

	DOMAIN := "neuron-nexus.ru"
	m := &autocert.Manager{
		Cache:      autocert.DirCache("../../var2/www/.cache"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(DOMAIN, "www."+DOMAIN, "doc."+DOMAIN, "docs."+DOMAIN),
		Email:      "info@realnoevremya.ru",
	}
	tlsServer := &http.Server{
		Addr: ":8000",
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		Handler: e.router,
	}

	log.Fatal(tlsServer.ListenAndServeTLS("", ""))
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
