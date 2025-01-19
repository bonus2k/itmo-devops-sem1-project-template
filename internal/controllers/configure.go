package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"project_sem/internal/logger"
	"project_sem/internal/services"
)

func MetricsRouter(logg *logger.Logger, ps *services.PriceService) chi.Router {
	ctrl := NewPriceHandler(logg, ps)
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RequestID,
		InjectLogger(logg),
		//rest.GzipReqDecompression,
		//rest.GzipResCompression,
	)

	router.Route("/api/v0", func(r chi.Router) {
		r.Post("/prices", ctrl.SavePrice)
		r.Get("/prices", ctrl.GetPrice)
	})
	return router
}

func InjectLogger(logg *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logg.ContextWithLogger(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
