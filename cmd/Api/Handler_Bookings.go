package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/internal/middleware"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func Helperrole(r *http.Request) (uuid.UUID, error) {

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
		Useridstr := chi.URLParam(r, "iduser")

		UserIDpars, erro := uuid.Parse(Useridstr)

		if erro != nil {
			return uuid.Nil, erro
		}

		UserID = UserIDpars
	}
	return UserID, nil
}

func (app Application) CreateBooking(w http.ResponseWriter, r *http.Request) {

	type param struct {
		status     string
		start_date time.Time
		end_dat    time.Time
	}
	decode := json.NewDecoder(r.Body)
	params := param{}
	err := decode.Decode(&params)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendecode param %v", err))
		return
	}
	id_petstr := chi.URLParam(r, "id_pet")

	if id_petstr == "" {
		jsonresponse.RespondWithNotfound(w, "gagal mendapatkan id pet dari url")
		return
	}

	id_pet, erro := uuid.Parse(id_petstr)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse id dari pet %v", erro))
		return
	}

	id_user, errr := Helperrole(r)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal authentikasi role %v", errr))
	}

	bookings, erros := app.Service.CreateNewBookings(r.Context(), id_pet, id_user, params.status, params.start_date, params.end_dat)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendaptakn bookings dari database %v", erros))
		return
	}

	jsonresponse.ResponSuccess(w, 200, bookings)
}
