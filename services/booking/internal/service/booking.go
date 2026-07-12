package service

type Repository interface{}

type Booking struct {
	repository Repository
}

func New(
	repository Repository,
) *Booking {
	return &Booking{
		repository: repository,
	}
}
