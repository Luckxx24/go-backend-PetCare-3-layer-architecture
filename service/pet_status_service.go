package service

import (
	"context"
	"database/sql"
	"errors"
	"pet-care/database"
	"pet-care/internal/middleware"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IsValidstatspsl(r string) bool {
	if r == "makan" || r == "sehat" || r == "sakit" || r == "grooming" {
		return true
	}
	return false
}

func (S *Services) CreatePetlOG(ctx context.Context, status, photo_path, note string, ID, id_booking uuid.UUID) (database.PetStatusLog, error) {
	role, okey := middleware.GetRoleFromContext(ctx)
	if !okey {
		return database.PetStatusLog{}, errors.New("role tidak ditemukan di dalam context")
	}

	oke := IsValidRole(role)

	if !oke {
		return database.PetStatusLog{}, errors.New("role tidak valid")
	}

	if role == "User" {
		return database.PetStatusLog{}, errors.New("anda tidak memiliki akses untuk ini")
	}
	status = strings.TrimSpace(status)
	note = strings.TrimSpace(note)

	if status == "" || note == "" {
		return database.PetStatusLog{}, errors.New("masukan status dan note")
	}

	ok := IsValidstatspsl(status)

	if !ok {
		return database.PetStatusLog{}, errors.New("enum status tidak valid")
	}

	bookings, err := S.StoreDB.Bookings.GetBookingByUserID(ctx, database.GetBookingByUserIDParams{
		ID:     id_booking,
		UserID: ID,
	})

	if bookings.Status == "PENDING" {
		return database.PetStatusLog{}, errors.New("mohon ubah status booking")
	}

	if err != nil {
		return database.PetStatusLog{}, err
	}

	notes := sql.NullString{
		String: note,
		Valid:  true,
	}
	photo := sql.NullString{
		String: photo_path,
		Valid:  true,
	}

	PetLog, err := S.StoreDB.Pet_Status_Log.CreateNewLog(ctx, database.CreateNewLogParams{
		ID:         uuid.New(),
		IDBookings: bookings.ID,
		Status:     database.Conds(status),
		Note:       notes,
		PhotoPath:  photo,
		CreatedBy:  ID,
		CreatedAt:  time.Now(),
		EditedAt:   time.Now(),
	})
	if err != nil {
		return database.PetStatusLog{}, err
	}
	return PetLog, nil
}

func (S *Services) GetAllLog(ctx context.Context, Page, PageSize int) ([]database.GetAllLogRow, error) {
	role, okey := middleware.GetRoleFromContext(ctx)
	if !okey {
		return []database.GetAllLogRow{}, errors.New("role tidak ditemukan di dalam context")
	}

	ok := IsValidRole(role)

	if !ok {
		return []database.GetAllLogRow{}, errors.New("masukan enum role yang valid")
	}

	if role == "User" {
		return []database.GetAllLogRow{}, errors.New("anda tidak memiliki hak ini")
	}

	UserIDstr, _ := middleware.GetIDFromContext(ctx)
	UserID, _ := uuid.Parse(UserIDstr)

	Offset := (Page - 1) * PageSize

	Alllog, err := S.StoreDB.Pet_Status_Log.GetAllLog(ctx, database.GetAllLogParams{
		ID:     UserID,
		Offset: int32(Offset),
		Limit:  int32(PageSize),
	})

	if err != nil {
		return []database.GetAllLogRow{}, err
	}

	return Alllog, nil

}

func (S *Services) GetlogUser(ctx context.Context, id_bookings uuid.UUID) (database.GetLogRow, error) {
	UsersIDstr, oke := middleware.GetIDFromContext(ctx)

	if !oke {
		return database.GetLogRow{}, errors.New("tidak bisa mendapatkan id dari context")
	}
	UserID, errr := uuid.Parse(UsersIDstr)

	if errr != nil {
		return database.GetLogRow{}, errr
	}

	bookings, erro := S.StoreDB.Bookings.GetBookingByUserID(ctx, database.GetBookingByUserIDParams{
		ID:     UserID,
		UserID: id_bookings,
	})

	if erro != nil {
		return database.GetLogRow{}, erro
	}

	LogUser, err := S.StoreDB.Pet_Status_Log.GetLog(ctx, bookings.ID)

	if err != nil {
		return database.GetLogRow{}, err
	}

	return LogUser, nil
}

func (S *Services) UpdateLog(ctx context.Context, ID, id_bookings, id_userp uuid.UUID, status, note, path string) (database.PetStatusLog, error) {
	role, okey := middleware.GetRoleFromContext(ctx)
	if !okey {
		return database.PetStatusLog{}, errors.New("role tidak ditemukan di dalam context")
	}

	oke := IsValidRole(role)

	if !oke {
		return database.PetStatusLog{}, errors.New("role tidak valid")
	}

	if role == "User" {
		return database.PetStatusLog{}, errors.New("anda tidak memiliki akses untuk ini")
	}
	status = strings.TrimSpace(status)
	note = strings.TrimSpace(note)

	if status == "" || note == "" {
		return database.PetStatusLog{}, errors.New("masukan status dan note")
	}

	booking, err := S.StoreDB.Bookings.GetBookingByUserID(ctx, database.GetBookingByUserIDParams{
		ID:     id_bookings,
		UserID: id_userp,
	})

	if err != nil {
		return database.PetStatusLog{}, err
	}

	if booking.Status == "PENDING" {
		return database.PetStatusLog{}, errors.New("status masih pending harap rubah status")
	}

	ok := IsValidstatspsl(status)

	if !ok {
		return database.PetStatusLog{}, errors.New("masukan enum status yang valid")
	}

	Notes := sql.NullString{
		String: note,
		Valid:  true,
	}

	path_photo := sql.NullString{
		String: path,
		Valid:  true,
	}

	UpdateLOg, errs := S.StoreDB.Pet_Status_Log.UpdateLog(ctx, database.UpdateLogParams{
		Status:    database.Conds(status),
		Note:      Notes,
		PhotoPath: path_photo,
		ID:        ID,
		EditedAt:  time.Now(),
	})

	if errs != nil {
		return database.PetStatusLog{}, errs
	}

	return UpdateLOg, nil
}

func (S *Services) DeleteLogPet(ctx context.Context, ID uuid.UUID) (bool, error) {
	role, _ := middleware.GetRoleFromContext(ctx)

	if ok := IsValidRole(role); !ok {
		return false, errors.New("masukan enum role yang valid")
	}

	if role == "Users" {
		return false, errors.New("anda tidak memiliki akses")
	}
	err := S.StoreDB.Pet_Status_Log.DeleteLog(ctx, ID)

	if err != nil {
		return false, err
	}

	return true, nil
}
