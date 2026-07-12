package router

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/handler"
)

func New(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

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

			// r.With(authMW).
			// 	Post("/{portId}/bookings", h.CreateBooking)

			r.Group(func(r chi.Router) {
				// r.Use(authMW)
				// r.Use(adminMW)

				r.Post("/", h.CreatePort)
				r.Patch("/{portId}", h.UpdatePort)
				r.Delete("/{portId}", h.DeletePort)
			})
		})

		r.Group(func(r chi.Router) {
			// r.Use(authMW)
			// r.Use(adminMW)

			r.Post("/", h.CreateStation)
			r.Put("/{id}", h.UpdateStation)
			r.Delete("/{id}", h.DeleteStation)
		})
	})

	r.Group(func(r chi.Router) {
		// r.Use(authMW)

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
