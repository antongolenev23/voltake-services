package router

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/config"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/handler"
	appmiddleware "github.com/antongolenev23/voltake-services/services/booking/internal/http-server/middleware"
)

func New(h *handler.Handler, cfg *config.ConfigJWT) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		appmiddleware.RequestID(),
		appmiddleware.LoggerMiddleware(h.Log),
	)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	r.Route("/stations", func(r chi.Router) {
		r.Get("/", h.GetStations)
		r.Get("/{id}", h.GetStation)

		r.Route("/{stationId}/ports", func(r chi.Router) {
			r.Get("/", h.GetPorts)
			r.Get("/{portId}", h.GetPort)

			r.With(appmiddleware.JWT(cfg.Secret)).
				Post("/{portId}/bookings", h.CreateBooking)

			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.JWT(cfg.Secret))
				// r.Use(adminMW)

				r.Post("/", h.CreatePort)
				r.Patch("/{portId}", h.UpdatePort)
				r.Delete("/{portId}", h.DeletePort)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.JWT(cfg.Secret))
			// r.Use(adminMW)

			r.Post("/", h.CreateStation)
			r.Put("/{id}", h.UpdateStation)
			r.Delete("/{id}", h.DeleteStation)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.JWT(cfg.Secret))

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", h.GetMe)
		})

		r.Route("/bookings", func(r chi.Router) {
			r.Get("/", h.GetBookings)
			r.Get("/{id}", h.GetBooking)
			r.Patch("/{id}/cancel", h.CancelBooking)
		})
	})

	return r
}
