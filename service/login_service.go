package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (S *Services) Login(ctx context.Context, nama, password, email string) error {
	nama = strings.TrimSpace(nama)
	password = strings.TrimSpace(password)
	email = strings.TrimSpace(email)

	if nama == "" || password == "" || email == "" {
		return errors.New("harap isi kolom login")
	}

	user, err := S.StoreDB.Users.GetUseremail(ctx, email)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("User tidak ditemukan")
		}
		return err
	}

	if nama != user.Nama {
		return errors.New("nama tidak ditemukan")
	}

	erro := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if erro != nil {
		return errors.New("gagal menghash passowrd")
	}

	return nil

}
