package main

import (
	"net/http"

	"github.com/canhlinh/cointracker/api"
	"github.com/canhlinh/cointracker/backend"
	chi "github.com/go-chi/chi/v5"
)

func main() {

	app := backend.NewApp()
	go app.Run()

	apihandler := &api.Handler{
		App:    app,
		Router: chi.NewRouter(),
	}
	api.Route(apihandler)
	http.ListenAndServe(":8000", apihandler.Router)
}
