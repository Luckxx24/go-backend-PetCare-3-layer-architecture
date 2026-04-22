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
	"golang.org/x/crypto/bcrypt"
)

func (S *Services) CreateUser(ctx context.Context, nama, password, email, role string) (database.User, error) {

	nama = strings.TrimSpace(nama)
	email = strings.TrimSpace(email)
	if nama == "" || email == "" || len(password) < 8 {
		return database.User{}, errors.New("pastikan kolom nama,email di isi dan password lebih dar 8 karakter ")
	}

	_, err := S.StoreDB.Users.GetUseremail(ctx, email)

	if err == nil {
		return database.User{}, errors.New("Email sudah terdaftar harap masukan email baru")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return database.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return database.User{}, errors.New("gagal mendecrypt password")
	}

	okey := IsValidRole(role)

	if !okey {
		return database.User{}, errors.New("Role tidak termasuk dalam kategor")
	}

	roleEnum := database.Role(role)

	waktusekarang := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	User, err := S.StoreDB.Users.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Nama:      nama,
		Password:  string(hash),
		Email:     email,
		Role:      roleEnum,
		CreatedAt: waktusekarang,
	})

	if err != nil {
		return database.User{}, err
	}

	return User, nil
}

func (S *Services) DeleteUser(ctx context.Context, ID uuid.UUID) (error, bool) {
	role, ok := middleware.GetRoleFromContext(ctx)
	if !ok {
		return errors.New("gagal mendapatkan role dari context"), false
	}

	okey := IsValidRole(role)

	if !okey {
		return errors.New("tidak dapat mendapatkan role dari context"), false
	}

	if role == "User" {
		return errors.New("Anda tida memiliki akses ini "), false
	}

	err := S.StoreDB.Users.DeleteUser(ctx, ID)

	if err != nil {
		return errors.New("Gagal menghapus data"), false
	}

	return nil, true

}

func (S *Services) UpdateUser(ctx context.Context, ID uuid.UUID, nama, password, email, role string) (database.User, error) {

	nama = strings.TrimSpace(nama)
	email = strings.TrimSpace(email)
	if nama == "" || email == "" || len(password) < 8 {
		return database.User{}, errors.New("Mohon masukan nama dan pastikan karakter password lebih dari 8")
	}

	_, err := S.StoreDB.Users.GetUseremail(ctx, email)

	if err == nil {
		return database.User{}, errors.New("masukan email yang belum terdaftar")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return database.User{}, err
	}

	passwords, erros := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if erros != nil {
		return database.User{}, erros
	}

	Users, err := S.StoreDB.Users.UpdateUser(ctx, database.UpdateUserParams{
		Nama:     nama,
		Password: string(passwords),
		Email:    email,
		Role:     database.Role(role),
		ID:       ID,
	})

	if err != nil {
		return database.User{}, err
	}

	return Users, nil

}

func IsValidRole(r string) bool {
	if r == "User" || r == "Staff" || r == "Admin" {
		return true
	}
	return false
}

func (S *Services) ListsUserID(ctx context.Context, Page, pagesize int) ([]database.ListsUserRow, error) {

	role, ok := middleware.GetRoleFromContext(ctx)
	if !ok {
		return []database.ListsUserRow{}, errors.New("gagal mendapatkan role dari context")
	}

	Okey := IsValidRole(role)

	if !Okey {
		return []database.ListsUserRow{}, errors.New("tidak bisa mengambil role dari context")
	}

	if role == "Users" {
		return []database.ListsUserRow{}, errors.New("unauthorized")
	}

	offset := (Page - 1) * pagesize

	ListUser, err := S.StoreDB.Users.ListsUserID(ctx, database.ListsUserParams{
		Limit:  int32(pagesize),
		Offset: int32(offset),
	})
	if err != nil {
		return []database.ListsUserRow{}, err
	}

	return ListUser, nil
}
