package api

import (
	chi "github.com/go-chi/chi/v5"
)

func Route(handler *Handler) {
	handler.Router.Route("/api/v1", func(r chi.Router) {
		r.Route("/configs", func(r chi.Router) {
			r.Get("/", handler.ApiHandler(GetAlertConfigs))
			r.Post("/", handler.ApiHandler(CreateAlertConfigs))
			r.Post("/import", handler.ApiHandler(ImportSymbols))
			r.Get("/{id}", handler.ApiHandler(GetAlertConfig))
			r.Put("/{id}", handler.ApiHandler(UpdateAlertConfig))
			r.Delete("/{id}", handler.ApiHandler(DeleteAlertConfig))
			r.Delete("/", handler.ApiHandler(DeleteAlertConfigs))
		})
		r.Route("/timeframes", func(r chi.Router) {
			r.Get("/", handler.ApiHandler(GetTimeframes))
			r.Get("/{id}", handler.ApiHandler(GetTimeframe))
			r.Put("/{id}", handler.ApiHandler(UpdateTimeframe))
		})
		r.Get("/klines", handler.ApiHandler(GetKlines))
		r.Get("/alerts", handler.ApiHandler(GetAlerts))
		r.Get("/symbols", handler.ApiHandler(GetSymbols))
	})
}
