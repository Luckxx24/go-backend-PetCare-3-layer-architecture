package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/database"
	"pet-care/internal/middleware"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

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
		Useridstr := chi.URLParam(r, "iduser")

		UserIDpars, erro := uuid.Parse(Useridstr)

		if erro != nil {
			return uuid.Nil, erro
		}

		UserID = UserIDpars
	}
	return UserID, nil
}

func HelperIDBookings(r *http.Request) (uuid.UUID, error) {
	BookIDstr := chi.URLParam(r, "id_booking")

	if BookIDstr == "" {
		return uuid.Nil, errors.New("Not Fund")
	}

	BooKID, erro := uuid.Parse(BookIDstr)

	if erro != nil {
		return uuid.Nil, erro
	}
	return BooKID, nil
}

type params struct {
	status     string
	start_date time.Time
	end_dat    time.Time
}

func (app Application) CreateBooking(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	param := params{}
	err := decode.Decode(&param)

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

	id_user, errr := HelperroleGetID(r)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal authentikasi role %v", errr))
	}

	bookings, erros := app.Service.CreateNewBookings(r.Context(), id_pet, id_user, param.status, param.start_date, param.end_dat)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendaptakn bookings dari database %v", erros))
		return
	}

	jsonresponse.ResponSuccess(w, 200, bookings)
}

func (app Application) GetBookingmany(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Status string
	}

	decode := json.NewDecoder(r.Body)
	params := param{}
	erro := decode.Decode(&params)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, "gagal mendecode parasm")
		return
	}
	pagestr := r.URL.Query().Get("page")
	pagesizestr := r.URL.Query().Get("pagesize")

	var page int
	var pagesize int

	page = 1
	pagesize = 10

	if pagestr != "" {
		p, erro := strconv.Atoi(pagestr)

		if erro != nil || p < 0 {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("error ketika men parse page %v", erro))
			return
		}

		page = p
	}

	if pagesizestr != "" {
		p, errr := strconv.Atoi(pagesizestr)

		if errr != nil || p < 0 {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("error ketika men parse pagesize %v", errr))
			return
		}

		pagesize = p
	}

	pets, err := app.Service.GetBookingByStatus(r.Context(), params.Status, page, pagesize)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan booking %v", err))
	}

	jsonresponse.ResponSuccess(w, 200, pets)
}

func (app Application) UpdateBookings(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	param := params{}
	erro := decode.Decode(&param)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendecode param %v", erro))
		return
	}
	IDuser, errr := HelperroleGetID(r)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pet id %v", errr))
		return
	}

	BookID, erros := HelperIDBookings(r)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatka ID Booking %v", erros))
	}

	booking, errs := app.Service.StoreDB.Bookings.GetBookingByUserID(r.Context(), database.GetBookingByUserIDParams{
		ID:     BookID,
		UserID: IDuser,
	})

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pets %v", errs))
		return
	}

	bookings, err := app.Service.UpdateBookings(r.Context(), param.status, booking.ID, param.start_date, param.end_dat)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mengudpate booking %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, bookings)
}

func (app Application) DeleteBooking(w http.ResponseWriter, r *http.Request) {

	ID, errs := HelperIDBookings(r)

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID Bookings %v", errs))
		return
	}
	err := app.Service.DeleteBooking(r.Context(), ID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menghapus data %v", err))
	}
}
