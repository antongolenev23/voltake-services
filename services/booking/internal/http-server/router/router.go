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
		r.Route("/{stationID}/ports", func(r chi.Router) {
			r.Get("/{portID}", h.GetPort)

			r.With(appmiddleware.Auth(cfg.Secret)).
				Post("/{portID}/bookings", h.CreateBooking)

			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.Auth(cfg.Secret), appmiddleware.IsAdmin)

				r.Post("/", h.CreatePort)
				r.Post("/{portID}/activate", h.ActivatePort)
				r.Post("/{portID}/deactivate", h.DeactivatePort)
				r.Delete("/{portID}", h.DeletePort)
			})
		})

		r.Get("/{stationID}", h.GetStation)
		r.Get("/", h.GetStations)

		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Auth(cfg.Secret), appmiddleware.IsAdmin)

			r.Post("/", h.CreateStation)
			r.Put("/{stationID}", h.UpdateStation)
			r.Delete("/{stationID}", h.DeleteStation)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.Auth(cfg.Secret))

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", h.GetMe)
		})

		r.Route("/bookings", func(r chi.Router) {
			r.Get("/", h.GetBookings)
			r.Get("/{bookingID}", h.GetBooking)
			r.Patch("/{bookingID}/cancel", h.CancelBooking)
		})
	})

	return r
}
