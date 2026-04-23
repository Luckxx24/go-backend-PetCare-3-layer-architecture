package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app Application) CreatePet(w http.ResponseWriter, r *http.Request) {

	type param struct {
		Nama  string
		Jenis string
		Age   int
	}

	decode := json.NewDecoder(r.Body)
	params := param{}

	erro := decode.Decode(&params)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendecode data %v", erro))
	}

	role, ok := middleware.GetRoleFromContext(r.Context())

	if !ok {
		jsonresponse.RespondWithBadRequest(w, "gagal mengambil role dari context")
	}

	var userID uuid.UUID

	if role == "User" {
		UsersIDparse, ok := middleware.GetIDFromContext(r.Context())
		if !ok {
			jsonresponse.RespondWithBadRequest(w, "gagal mengambil ID dari context")
		}
		userIDI, errs := uuid.Parse(UsersIDparse)

		if errs != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse id dari context %v", errs))
		}

		userID = userIDI
	} else {
		UserIDstr := chi.URLParam(r, "iduser")

		userIDI, errr := uuid.Parse(UserIDstr)

		if errr != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse id dari context %v", errr))
		}

		userID = userIDI
	}
	pet, err := app.Service.CreatePets(r.Context(),
		params.Nama,
		params.Jenis,
		params.Age,
		userID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pet %v", err))
	}

	jsonresponse.ResponSuccess(w, 200, pet)

}

func (app Application) DeletePets(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.GetRoleFromContext(r.Context())

	if !ok {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan role dari context")
	}

	var UserID uuid.UUID

	if role == "User" {
		Useridstr, ok := middleware.GetIDFromContext(r.Context())
		useridpars, errr := uuid.Parse(Useridstr)

		if errr != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal parse ID %v", errr))
		}

		if !ok {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari context")
		}
		UserID = useridpars
	}

	if role == "Staff" || role == "Admin" {
		Useridstr := chi.URLParam(r, "iduser")

		UserIDpars, erro := uuid.Parse(Useridstr)

		if erro != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal parse ID %v", erro))
		}

		UserID = UserIDpars
	}

	err := app.Service.DeletePets(r.Context(), UserID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menghapus data %v", err))
	}
}
