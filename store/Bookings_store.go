package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Bookings interface {
	CreateNewBookings(ctx context.Context, arg database.CreateNewBookingsParams) (database.Booking, error)
	DeleteBooking(ctx context.Context, id uuid.UUID) error
	GetBooking(ctx context.Context, id uuid.UUID) (database.GetBookingRow, error)
	GetBookingByStatus(ctx context.Context, arg database.GetBookingByStatusParams) ([]database.GetBookingByStatusRow, error)
	UpdateBookings(ctx context.Context, arg database.UpdateBookingsParams) (database.Booking, error)
	GetBookingByUserID(ctx context.Context, arg database.GetBookingByUserIDParams) (database.GetBookingByUserIDRow, error)
}

type bookings struct {
	q *database.Queries
}

func (b bookings) CreateNewBookings(ctx context.Context, arg database.CreateNewBookingsParams) (database.Booking, error) {
	return b.q.CreateNewBookings(ctx, arg)
}

func (b bookings) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	return b.q.DeleteBooking(ctx, id)
}

func (b bookings) GetBooking(ctx context.Context, id uuid.UUID) (database.GetBookingRow, error) {
	return b.q.GetBooking(ctx, id)
}

func (b bookings) GetBookingByStatus(ctx context.Context, arg database.GetBookingByStatusParams) ([]database.GetBookingByStatusRow, error) {
	return b.q.GetBookingByStatus(ctx, arg)
}

func (b bookings) UpdateBookings(ctx context.Context, arg database.UpdateBookingsParams) (database.Booking, error) {
	return b.q.UpdateBookings(ctx, arg)
}

func (b bookings) GetBookingByUserID(ctx context.Context, arg database.GetBookingByUserIDParams) (database.GetBookingByUserIDRow, error) {
	return b.q.GetBookingByUserID(ctx, arg)
}
