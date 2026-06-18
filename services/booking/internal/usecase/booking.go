package usecase

type Storage interface{}

type Booking struct {
	storage Storage
}

func New(
	storage Storage,
) *Booking {
	return &Booking{
		storage: storage,
	}
}
