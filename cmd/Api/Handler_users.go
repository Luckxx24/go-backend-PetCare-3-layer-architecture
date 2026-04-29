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

type paramuser struct {
	Nama     string
	Password string
	Email    string
	Role     string
}

func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := paramuser{}

	err := decode.Decode(&params)

	if err != nil {
		jsonresponse.RespondError(w, 401, fmt.Sprintf("error ketika decode,%v", err))
		return
	}

	User, erro := app.Service.CreateUser(r.Context(), params.Nama, params.Password, params.Email, params.Role)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("error ketika mengupload data user %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, User)

}

func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
	UserIDstr, ok := middleware.GetIDFromContext(r.Context())

	if !ok {
		jsonresponse.RespondWithBadRequest(w, ("gagal mendapatkan ID dari context"))
		return
	}

	UserID, err := uuid.Parse(UserIDstr)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal parse ID %v", err))
		return
	}

	User, erro := app.Service.StoreDB.Users.GetUserID(r.Context(), UserID)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan user dari database %v", erro))
		return
	}

	jsonresponse.ResponSuccess(w, 200, User)
}

func (app *Application) UpdateUser(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := paramuser{}
	err := decode.Decode(&params)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, "gagal men-decode params")
		return
	}

	UserIDreal, errr := HelperroleGetID(r)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID %v", errr))
		return
	}

	UpdatedUser, erros := app.Service.UpdateUser(r.Context(), UserIDreal, params.Nama, params.Password, params.Email, params.Role)

	if erros != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal meng-update User %v", erros))
		return
	}

	jsonresponse.ResponSuccess(w, 200, UpdatedUser)
}

func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {

	UserIDstr := chi.URLParam(r, "id_user")

	UserID, erro := uuid.Parse(UserIDstr)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse %v", erro))
		return
	}

	err := app.Service.DeleteUser(r.Context(), UserID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal meng hapus data %v", err))
		return
	}
}

func (app *Application) GetListUser(w http.ResponseWriter, r *http.Request) {

	Page, PageSize, err := HelperPage(r)

	if err != nil {
		jsonresponse.RespondWithNotfound(w, fmt.Sprintf("gagal mendapatkan data page %v", err))
		return
	}

	User, erro := app.Service.ListsUserID(r.Context(), Page, PageSize)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan data %v", erro))
		return
	}

	jsonresponse.ResponSuccess(w, 200, User)

}
