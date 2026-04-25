package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/database"
	"pet-care/internal/middleware"
	"strconv"

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

	role, ok := middleware.GetRoleFromContext(r.Context())

	if !ok {
		jsonresponse.RespondWithBadRequest(w, "gagal mengambil role dari context")
		return
	}

	var userID uuid.UUID

	if role == "User" {
		UsersIDparse, ok := middleware.GetIDFromContext(r.Context())
		if !ok {
			jsonresponse.RespondWithBadRequest(w, "gagal mengambil ID dari context")
			return
		}
		userIDI, errs := uuid.Parse(UsersIDparse)

		if errs != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse id dari context %v", errs))
			return
		}

		userID = userIDI
	} else {
		UserIDstr := chi.URLParam(r, "iduser")

		userIDI, errr := uuid.Parse(UserIDstr)

		if errr != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse id dari context %v", errr))
			return
		}

		userID = userIDI
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
	return

}

func (app Application) UpdatePets(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := param{}

	erro := decode.Decode(&params)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendecode data %v", erro))
	}
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
			return
		}

		if !ok {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari context")
			return
		}
		UserID = useridpars
	}

	if role == "Staff" || role == "Admin" {
		Useridstr := chi.URLParam(r, "iduser")

		UserIDpars, erro := uuid.Parse(Useridstr)

		if erro != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal parse ID %v", erro))
			return
		}

		UserID = UserIDpars
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

	pet, err := app.Service.GetPetUser(r.Context(), page, pagesize)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("error ketika mendapatkan pet %v ", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, pet)
	return

}

func (app Application) GetPetsadmin(w http.ResponseWriter, r *http.Request) {

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

	pets, erro := app.Service.GetPetsStaff(r.Context(), page, pagesize)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan pets dari database %v", erro))
		return
	}
	jsonresponse.ResponSuccess(w, 200, pets)
	return

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
	return

}

func (app *Application) DeletePets(w http.ResponseWriter, r *http.Request) {
	role, oke := middleware.GetRoleFromContext(r.Context())

	if !oke {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan role dari context")
		return
	}

	var UserID uuid.UUID

	if role == "Staff" || role == "Admin" {
		idusertstr := chi.URLParam(r, "user_id")

		if idusertstr == "" {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan id dari url param")
			return
		}

		id, err := uuid.Parse(idusertstr)

		if err != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men-parse UserID %v", err))
			return
		}

		UserID = id

	}

	if role == "User" {
		IDuserstr, okey := middleware.GetIDFromContext(r.Context())

		if !okey {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan id dari context"))
			return
		}

		id, errors := uuid.Parse(IDuserstr)

		if errors != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men-parse id %v", errors))
			return
		}

		UserID = id

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
