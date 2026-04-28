package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func HelperIDPSL(r *http.Request) (uuid.UUID, error) {
	PSLIDstr := chi.URLParam(r, "id_petstatuslog")

	if PSLIDstr == "" {
		return uuid.Nil, errors.New("Not Fund")
	}

	PslID, erro := uuid.Parse(PSLIDstr)

	if erro != nil {
		return uuid.Nil, erro
	}
	return PslID, nil
}

type parampet struct {
	status     string
	photo_path string
	note       string
}

func (app Application) CreatePetStatusLOG(W http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := parampet{}
	erros := decode.Decode(&params)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(W, fmt.Sprintf("gagal men decode %v", erros))
	}

	UserID, Erro := HelperroleGetID(r)

	if Erro != nil {
		jsonresponse.RespondWithBadRequest(W, fmt.Sprintf("gagal menagmbil ID User %v"))
		return
	}
	IDBooking, erro := HelperIDBookings(r)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(W, fmt.Sprintf("gagal mendapatka ID Booking %v", erro))
	}

	PetLOg, err := app.Service.CreatePetlOG(r.Context(), params.status, params.photo_path, params.note, UserID, IDBooking)

	if err != nil {
		jsonresponse.RespondWithBadRequest(W, fmt.Sprintf("gagal menyimpan PetStatusLOg %v", err))
		return
	}

	jsonresponse.ResponSuccess(W, 200, PetLOg)
}

func (app Application) GetPetlog(w http.ResponseWriter, r *http.Request) {
	page, pagesize, erros := HelperPage(r)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan info page dari url %v", erros))
	}
	petlogmany, err := app.Service.GetAllPetLog(r.Context(), page, pagesize)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan log pet %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, petlogmany)
}

func (app Application) GetPetLOgUser(w http.ResponseWriter, r *http.Request) {
	BookingID, errs := HelperIDBookings(r)

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID booking %v", errs))
		return
	}
	PetLOgUser, err := app.Service.GetpetlogUser(r.Context(), BookingID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menemukan Pet LOg User %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, PetLOgUser)
}

func (app Application) UpdatePetLOg(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := parampet{}
	erros := decode.Decode(&params)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men decode %v", erros))
	}

	UserID, errs := HelperroleGetID(r)
	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID User %v", errs))
		return
	}
	PslID, errr := HelperIDPSL(r)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID PSL %v", errr))
		return
	}

	BookingID, erro := HelperIDBookings(r)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID PSL %v", errr))
		return
	}

	UPdatedPet, err := app.Service.UpdateLogpet(r.Context(), PslID, BookingID, UserID, params.status, params.note, params.photo_path)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal meng-updated pet %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, UPdatedPet)

}

func (app Application) DeletePetLOg(w http.ResponseWriter, r *http.Request) {

	PslID, erro := HelperIDPSL(r)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID PSL %v", erro))
		return
	}
	err := app.Service.DeleteLogPet(r.Context(), PslID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menghapus PETS %v", err))
		return
	}
}
