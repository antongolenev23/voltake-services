package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/handler"
)

func New(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Auth
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", h.Register)
		auth.Post("/login", h.Login)
	})

	// Users
	r.Route("/users", func(users chi.Router) {
		users.Get("/me", h.GetMe)
	})

	// Stations
	r.Route("/stations", func(stations chi.Router) {
		stations.Get("/", h.GetStations)
		stations.Post("/", h.CreateStation) // admin
		stations.Get("/{id}", h.GetStation)
		stations.Put("/{id}", h.UpdateStation)    // admin
		stations.Delete("/{id}", h.DeleteStation) // admin

		stations.Route("/{stationId}/ports", func(ports chi.Router) {
			ports.Get("/", h.GetPorts)
			ports.Post("/", h.CreatePort) // admin
			ports.Get("/{portId}", h.GetPort)
			ports.Patch("/{portId}", h.UpdatePort)  // admin
			ports.Delete("/{portId}", h.DeletePort) // admin

			ports.Post("/{portId}/bookings", h.CreateBooking)
		})
	})

	// Bookings
	r.Route("/bookings", func(bookings chi.Router) {
		bookings.Get("/", h.GetBookings)
		bookings.Get("/{id}", h.GetBooking)
		bookings.Patch("/{id}/cancel", h.CancelBooking)
	})

	return r
}
