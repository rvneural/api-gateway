package main

import "rvneural/api-gateway/internal/pkg/app"

func main() {
	a := app.New()
	a.Start()
}
