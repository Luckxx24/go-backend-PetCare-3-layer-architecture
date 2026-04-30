package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/database"
	"pet-care/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type param struct {
	Nama         string
	Jenis        string
	Age          int
	Catatan      string
	Berat        string
	JenisKelamin string
	Ras          string
	IsVaxinated  bool
	PhotoPath    string
}

func (app Application) CreatePet(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := param{}

	erro := decode.Decode(&params)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendecode data %v", erro))
		return
	}

	userID, errs := HelperroleGetID(r)
	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID dari user %v", errs))
		return
	}

	pet, err := app.Service.CreatePets(r.Context(),
		params.Nama,
		params.Jenis,
		params.Catatan,
		params.Ras,
		params.PhotoPath,
		params.Berat,
		params.JenisKelamin,
		params.Age,
		userID,
		params.IsVaxinated)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pet %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, pet)

}

func (app Application) UpdatePets(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := param{}

	erro := decode.Decode(&params)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendecode data %v", erro))
	}

	UserID, errs := HelperroleGetID(r)
	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID dari user %v", errs))
		return
	}
	idpetsstr := chi.URLParam(r, "id_pet")

	if idpetsstr == "" {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari urlparam")
		return
	}
	PetID, erros := uuid.Parse(idpetsstr)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men-parse ID %v", erros))
		return
	}

	pets, err := app.Service.UpdatePets(r.Context(), params.Nama,
		params.Jenis,
		params.Catatan,
		params.Ras,
		params.PhotoPath,
		params.Berat,
		params.JenisKelamin,
		params.Age,
		PetID,
		UserID,
		params.IsVaxinated)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mengupdate data %v", err))
		return
	}
	jsonresponse.ResponSuccess(w, 200, pets)
}

func (app Application) GetPetslistUser(w http.ResponseWriter, r *http.Request) {

	page, pagesize, eros := HelperPage(r)

	if eros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan nomor page %v", eros))
		return
	}

	pet, err := app.Service.GetPetUser(r.Context(), page, pagesize)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("error ketika mendapatkan pet %v ", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, pet)

}

func (app Application) GetPetsadmin(w http.ResponseWriter, r *http.Request) {

	page, pagesize, erros := HelperPage(r)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan nomor page %v", erros))
		return
	}

	pets, erro := app.Service.GetPetsStaff(r.Context(), page, pagesize)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pets dari database %v", erro))
		return
	}
	jsonresponse.ResponSuccess(w, 200, pets)

}

func (app *Application) GetPetsDetail(w http.ResponseWriter, r *http.Request) {
	role, oke := middleware.GetRoleFromContext(r.Context())

	if !oke {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan role dari context")
		return
	}

	var petsid uuid.UUID

	if role == "Staff" || role == "Admin" {
		idpetstr := chi.URLParam(r, "pet_id")

		if idpetstr == "" {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari url param")
			return
		}

		p, err := uuid.Parse(idpetstr)

		if err != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men-parse UserID %v", err))
			return
		}

		petsid = p

	}

	if role == "User" {
		idpetstr := chi.URLParam(r, "pet_id")

		if idpetstr == "" {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari url param")
			return
		}

		idpet, err := uuid.Parse(idpetstr)

		idusrstr, okey := middleware.GetIDFromContext(r.Context())

		if !okey {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari context")
			return
		}

		id_user, erro := uuid.Parse(idusrstr)

		if erro != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse id %v", erro))
		}

		p, err := app.Store.Pets.GetPetsByID(r.Context(), database.GetPetsByIDParams{

			ID:     idpet,
			UserID: id_user,
		})

		if err != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatakn pets id dari database %v", err))
			return
		}

		petsid = p

	}

	pets, erro := app.Service.StoreDB.Pets.GetPetsDetail(r.Context(), petsid)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pets %v", erro))
		return
	}

	jsonresponse.ResponSuccess(w, 200, pets)

}

func (app *Application) DeletePets(w http.ResponseWriter, r *http.Request) {
	UserID, errs := HelperroleGetID(r)

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan id user %v", errs))
		return
	}

	petsidstr := chi.URLParam(r, "pets_id")

	if petsidstr == "" {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id pet dari url")
		return
	}

	PetID, erro := uuid.Parse(petsidstr)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men parse id %v", erro))
		return
	}
	err := app.Service.DeletePets(r.Context(), UserID, PetID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menghapus data %v", err))
		return
	}
}
