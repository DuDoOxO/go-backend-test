package main

import (
	"sync"

	"github.com/go-backend-test/service"
)

type App struct {
	restful service.Restful
}

func NewApp() *App {
	app := &App{}
	app.restful = *service.NewRestful()
	return app
}

func (a *App) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	a.restful.Run(&wg)
	wg.Wait()
}

func main() {
	app := NewApp()
	app.Run()
}
