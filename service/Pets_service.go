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

func (S *Services) CreatePets(ctx context.Context, nama, jenis, catatan, ras, photo_path, Berat, JenisKelamin string, age int, usersID uuid.UUID, isvaxinated bool) (database.Pet, error) {
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
		ID:           uuid.New(),
		UserID:       usersID,
		Nama:         nama,
		Jenis:        jenis,
		Age:          ages,
		CreatedAt:    timenow,
		Catatan:      catatan,
		Berat:        Berat,
		JenisKelamin: database.Kelamin(JenisKelamin),
		Ras:          ras,
		IsVaxinated:  isvaxinated,
		PhotoPath:    photo_path,
	})

	if err != nil {
		return database.Pet{}, err
	}

	return pets, nil
}

func (S *Services) DeletePets(ctx context.Context, UserID, ID uuid.UUID) error {

	err := S.StoreDB.Pets.DeletePets(ctx, database.DeletePetsParams{
		ID:     ID,
		UserID: UserID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (S *Services) GetPetsStaff(ctx context.Context, Page, PageSize int) ([]database.GetPetsListStRow, error) {

	role, okey := middleware.GetRoleFromContext(ctx)

	if !okey {
		return []database.GetPetsListStRow{}, errors.New("tidal bisa mendapatkan role dari context")
	}

	ok := IsValidRole(role)

	if !ok {
		return []database.GetPetsListStRow{}, errors.New("role tidak valid")
	}

	if role == "User" {
		return []database.GetPetsListStRow{}, errors.New("anda tidak memiliki akses")
	}
	Offset := (Page - 1) * PageSize

	PetsMany, err := S.StoreDB.Pets.GetPetsListSt(ctx, database.GetPetsListStParams{
		Offset: int32(Offset),
		Limit:  int32(PageSize),
	})

	if err != nil {
		return []database.GetPetsListStRow{}, err
	}
	return PetsMany, nil
}

func (S *Services) GetPetUser(ctx context.Context, Page, PageSize int) ([]database.GetPetsListUserRow, error) {
	userIDstr, ok := middleware.GetIDFromContext(ctx)

	if !ok {
		return []database.GetPetsListUserRow{}, errors.New("gagal mendapatkan ID dari context")
	}

	UserID, err := uuid.Parse(userIDstr)

	if err != nil {
		return []database.GetPetsListUserRow{}, err
	}

	Offset := (Page - 1) * PageSize

	petsmanyuser, errr := S.StoreDB.Pets.GetPetsListUser(ctx, database.GetPetsListUserParams{
		UserID: UserID,
		Offset: int32(Offset),
		Limit:  int32(PageSize),
	})

	if errr != nil {
		return []database.GetPetsListUserRow{}, errr
	}

	return petsmanyuser, nil
}

func (S *Services) UpdatePets(ctx context.Context, nama, jenis, catatan, ras, photo_path, Berat, JenisKelamin string, age int, ID, usersID uuid.UUID, isvaxinated bool) (database.Pet, error) {
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
		ID:           ID,
		UserID:       usersID,
		Nama:         nama,
		Jenis:        jenis,
		Age:          ages,
		Catatan:      catatan,
		Berat:        Berat,
		JenisKelamin: database.Kelamin(JenisKelamin),
		Ras:          ras,
		IsVaxinated:  isvaxinated,
		PhotoPath:    photo_path,
	})

	if err != nil {
		return database.Pet{}, err
	}

	return pets, nil
}
