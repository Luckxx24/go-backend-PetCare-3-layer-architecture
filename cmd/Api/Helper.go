package main

import (
	"errors"
	"net/http"
	"pet-care/internal/middleware"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func HelperPage(r *http.Request) (int, int, error) {
	pagestr := r.URL.Query().Get("page")
	pagesizestr := r.URL.Query().Get("pagesize")

	Page := 1
	PageSize := 10

	if pagestr != "" {
		p, err := strconv.Atoi(pagestr)

		if err != nil && Page > 0 {
			Page = p
		}

		return 0, 0, err
	}

	if pagesizestr != "" {
		ps, err := strconv.Atoi(pagesizestr)

		if err != nil && PageSize > 0 {
			PageSize = ps
		}

		return 0, 0, err
	}
	return Page, PageSize, nil
}

func HelperroleGetID(r *http.Request) (uuid.UUID, error) {

	role, ok := middleware.GetRoleFromContext(r.Context())

	if !ok {
		return uuid.Nil, errors.New("gagal mendapatkan role dari context")
	}

	var UserID uuid.UUID

	if role == "User" {
		Useridstr, ok := middleware.GetIDFromContext(r.Context())
		useridpars, errr := uuid.Parse(Useridstr)

		if errr != nil {

			return uuid.Nil, errr
		}

		if !ok {
			return uuid.Nil, errors.New("gagal mendapatkan id dari context")
		}
		UserID = useridpars
	}

	if role == "Staff" || role == "Admin" {
		Useridstr := chi.URLParam(r, "id_user")

		UserIDpars, erro := uuid.Parse(Useridstr)

		if erro != nil {
			return uuid.Nil, erro
		}

		UserID = UserIDpars
	}
	return UserID, nil
}
