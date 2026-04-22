package Store

import (
	"pet-care/database"
)

type Storage struct {
	Users          Users
	Pets           Pets
	Bookings       Bookings
	Pet_Status_Log Pet_Status_Log
	Message        Message
	Notifications  Notifications
}

func NewStorage(DB database.DBTX) Storage {

	q := database.New(DB)

	return Storage{
		Users:          &users{q: q},
		Pets:           &pets{q: q},
		Bookings:       &bookings{q: q},
		Pet_Status_Log: &pet_status_log{q: q},
		Message:        &message{q: q},
		Notifications:  &notifications{q: q},
	}
}
