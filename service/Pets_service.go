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

func (S *Services) CreatePets(ctx context.Context, nama, jenis string, age int, usersID uuid.UUID) (database.Pet, error) {
	nama = strings.TrimSpace(nama)
	jenis = strings.TrimSpace(jenis)

	if nama == "" || jenis == "" || age <= 0 {
		return database.Pet{}, errors.New("harap isi kolom nama,jenis,age")

	}

	ages := sql.NullInt32{
		Int32: int32(age),
		Valid: true,
	}

	timenow := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	pets, err := S.StoreDB.Pets.CreatePets(ctx, database.CreatePetsParams{
		ID:        uuid.New(),
		UserID:    usersID,
		Nama:      nama,
		Jenis:     jenis,
		Age:       ages,
		CreatedAt: timenow,
	})

	if err != nil {
		return database.Pet{}, err
	}

	return pets, nil
}

func (S *Services) DeletePets(ctx context.Context, ID uuid.UUID) error {
	IDuserstr, okey := middleware.GetIDFromContext(ctx)

	if !okey {
		errors.New("tidal bisa mendapatkan ID dari context")
	}
	UserID, errs := uuid.Parse(IDuserstr)

	if errs != nil {
		return errs
	}

	err := S.StoreDB.Pets.DeletePets(ctx, database.DeletePetsParams{
		ID:     ID,
		UserID: UserID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (S *Services) GetPetsMany(ctx context.Context, Page, PageSize int) ([]database.GetPetsManyRow, error) {

	role, okey := middleware.GetRoleFromContext(ctx)

	if !okey {
		return []database.GetPetsManyRow{}, errors.New("tidal bisa mendapatkan role dari context")
	}

	ok := IsValidRole(role)

	if !ok {
		return []database.GetPetsManyRow{}, errors.New("role tidak valid")
	}

	if role == "Users" {
		return []database.GetPetsManyRow{}, errors.New("anda tidak memiliki akses")
	}
	Offset := (Page - 1) * PageSize

	PetsMany, err := S.StoreDB.Pets.GetPetsMany(ctx, database.GetPetsManyParams{
		Offset: int32(Offset),
		Limit:  int32(PageSize),
	})

	if err != nil {
		return []database.GetPetsManyRow{}, err
	}

	return PetsMany, nil
}

func (S *Services) UpdatePets(ctx context.Context, nama, jenis string, age int, ID uuid.UUID) (database.Pet, error) {
	nama = strings.TrimSpace(nama)
	jenis = strings.TrimSpace(jenis)

	if nama == "" || jenis == "" || age <= 0 {
		return database.Pet{}, errors.New("harap isi kolom nama,jenis,age")

	}

	UsersIDparse, ok := middleware.GetIDFromContext(ctx)
	if !ok {
		return database.Pet{}, errors.New("gagal mendapatkan ID user dari context")
	}
	usersID, err := uuid.Parse(UsersIDparse)

	if err != nil {
		return database.Pet{}, err
	}

	ages := sql.NullInt32{
		Int32: int32(age),
		Valid: true,
	}

	pets, err := S.StoreDB.Pets.UpdatePets(ctx, database.UpdatePetsParams{
		Nama:   nama,
		Jenis:  jenis,
		Age:    ages,
		ID:     ID,
		UserID: usersID,
	})

	if err != nil {
		return database.Pet{}, err
	}

	return pets, nil
}
