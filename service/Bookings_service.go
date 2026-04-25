package service

import (
	"context"
	"database/sql"
	"errors"
	"pet-care/database"
	"pet-care/internal/middleware"
	"time"

	"github.com/google/uuid"
)

func IsValidstats(status string) bool {
	if status == "PENDING" || status == "APPROVED" || status == "DONE" {
		return true
	}
	return false
}

func (S *Services) CreateNewBookings(ctx context.Context, pet_id, user_id uuid.UUID, status string, startDate, endDate time.Time) (database.Booking, error) {

	okey := IsValidstats(status)
	if !okey {
		return database.Booking{}, errors.New("harap masukan status sesuai enum")
	}

	waktu := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if startDate.After(endDate) {
		return database.Booking{}, errors.New("start_date tidak boleh setelah end date")
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return database.Booking{}, errors.New("tanggal booking tidak boleh di masa lalau")
	}

	Bookings, err := S.StoreDB.Bookings.CreateNewBookings(ctx,
		database.CreateNewBookingsParams{
			ID:        uuid.New(),
			PetID:     pet_id,
			UserID:    user_id,
			StartDate: startDate,
			EndDate:   endDate,
			Status:    database.BookingStatus(status),
			CreatedAt: waktu,
		})

	if err != nil {
		return database.Booking{}, err
	}

	return Bookings, err
}

func (S *Services) GetBookingByStatus(ctx context.Context, status string, Page, PageSize int) ([]database.GetBookingByStatusRow, error) {
	role, ok := middleware.GetRoleFromContext(ctx)

	oke := IsValidRole(role)

	if !oke {
		return []database.GetBookingByStatusRow{}, errors.New("role tidak valid masukan role valid")
	}
	if !ok {
		return []database.GetBookingByStatusRow{}, errors.New("role tidak ditemukan di dalam context")
	}

	if role == "User" {
		return []database.GetBookingByStatusRow{}, errors.New("Anda tidak memiliki akses ini")
	}

	okey := IsValidstats(status)
	if !okey {
		return []database.GetBookingByStatusRow{}, errors.New("masukan status sesuai enum ")
	}

	Offset := (Page - 1) * PageSize

	Bookings, err := S.StoreDB.Bookings.GetBookingByStatus(ctx, database.GetBookingByStatusParams{
		Status: database.BookingStatus(status),
		Offset: int32(Offset),
		Limit:  int32(PageSize),
	})

	if err != nil {
		return []database.GetBookingByStatusRow{}, err
	}

	return Bookings, nil
}

func (S *Services) DeleteBooking(ctx context.Context, ID uuid.UUID) (bool, error) {
	role, ok := middleware.GetRoleFromContext(ctx)

	okey := IsValidRole(role)

	if !okey {
		return false, errors.New("role tidak valid masukan role valid")
	}

	if !ok {
		return false, errors.New("role tidak ditemukan di dalam context")
	}

	if role == "User" {
		return false, errors.New("anda tidak memiliki akses ini")
	}

	err := S.StoreDB.Bookings.DeleteBooking(ctx, ID)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (S *Services) UpdateBookings(ctx context.Context, status string, ID uuid.UUID, startDate, endDate time.Time) (database.Booking, error) {
	role, ok := middleware.GetRoleFromContext(ctx)

	if !ok {
		return database.Booking{}, errors.New("role tidak ditemukan di dalam context")
	}
	oke := IsValidRole(role)

	if !oke {
		return database.Booking{}, errors.New("role tidak valid masukan role valid")
	}
	if role == "User" {
		return database.Booking{}, errors.New("anda tidak memiliki akses ini")
	}

	okey := IsValidstats(status)

	if !okey {
		return database.Booking{}, errors.New("anda tidak memiliki hak akses ini")
	}

	if startDate.After(endDate) {
		return database.Booking{}, errors.New("start date tidak boleh setelah enddate")
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return database.Booking{}, errors.New("start date tidak boleh di masa lalu")
	}

	UpdateBookings, err := S.StoreDB.Bookings.UpdateBookings(ctx, database.UpdateBookingsParams{
		StartDate: startDate,
		EndDate:   endDate,
		Status:    database.BookingStatus(status),
		ID:        ID,
	})

	if err != nil {
		return database.Booking{}, err
	}

	return UpdateBookings, nil
}
