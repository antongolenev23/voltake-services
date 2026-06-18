package handler

import "net/http"

type Booking interface {
}

type Handler struct {
	booking Booking
}

func New(
	booking Booking,
) *Handler {
	return &Handler{
		booking: booking,
	}
}

//  Auth

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Users

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Stations

func (h *Handler) GetStations(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) CreateStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) UpdateStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) DeleteStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Ports

func (h *Handler) GetPorts(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetPort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) CreatePort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) UpdatePort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) DeletePort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Bookings

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	// TODO
}
