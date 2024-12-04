package app

import (
	"os"
	endpoint "rvneural/api-gateway/internal/endpoint/app"
)

type App struct {
	Text2TextURL  string
	Audio2TextURL string
	Text2ImageURL string
	RSSUrl        string
	DBURL         string
	MEDIAURL      string
	BgURL         string
	API_KEY       string
}

func New() *App {
	return &App{
		Text2TextURL:  os.Getenv("Text2TextURL"),
		Audio2TextURL: os.Getenv("Audio2TextURL"),
		Text2ImageURL: os.Getenv("Text2ImageURL"),
		RSSUrl:        os.Getenv("RSSURL"),
		DBURL:         os.Getenv("DBURL"),
		MEDIAURL:      os.Getenv("MEDIAURL"),
		BgURL:         os.Getenv("BGURL"),
		API_KEY:       os.Getenv("API_KEY"),
	}
}

func (a *App) Start() {
	endpointApp := endpoint.New()

	endpointApp.AddAPIEndpoint("/text/*path", a.Text2TextURL)
	endpointApp.AddAPIEndpoint("/audio/*path", a.Audio2TextURL)
	endpointApp.AddAPIEndpoint("/image/*path", a.Text2ImageURL)
	endpointApp.AddAPIEndpoint("/db/*path", a.DBURL)
	endpointApp.AddAPIEndpoint("/rss/*path", a.RSSUrl)
	endpointApp.AddAPIEndpoint("/media/*path", a.MEDIAURL)
	endpointApp.AddAPIEndpoint("/bg/*path", a.BgURL)

	endpointApp.SetApiKey(a.API_KEY)

	endpointApp.Start()
}
