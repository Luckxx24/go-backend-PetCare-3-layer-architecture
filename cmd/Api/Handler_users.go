package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/internal/middleware"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Nama     string
		Password string
		Email    string
		Role     string
	}

	decode := json.NewDecoder(r.Body)
	params := param{}

	err := decode.Decode(&params)

	if err != nil {
		jsonresponse.RespondError(w, 401, fmt.Sprintf("error ketika decode,%v", err))
	}

	User, erro := app.Service.CreateUser(r.Context(), params.Nama, params.Password, params.Email, params.Role)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("error ketika memgupload data user %v", err))
	}

	jsonresponse.ResponSuccess(w, 200, User)

}

func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
	UserIDstr, ok := middleware.GetIDFromContext(r.Context())

	if !ok {
		jsonresponse.RespondWithBadRequest(w, ("gagal mendapatkan ID dari context"))
	}

	UserID, err := uuid.Parse(UserIDstr)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal parse ID %v", err))
	}

	User, erro := app.Service.StoreDB.Users.GetUserID(r.Context(), UserID)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan user dari database %v", erro))
	}

	jsonresponse.ResponSuccess(w, 200, User)
}

func (app *Application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Nama     string
		Password string
		Email    string
		Role     string
	}
	decode := json.NewDecoder(r.Body)
	params := param{}
	err := decode.Decode(params)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men-decode params"))
	}

	role, ok := middleware.GetRoleFromContext(r.Context())
	UserIDstri, oke := middleware.GetIDFromContext(r.Context())

	if !oke {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan ID dari context")
	}

	if !ok {
		jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan role")
	}

	var UserIDreal uuid.UUID
	var errParse error

	if role == "Staff" || role == "Admin" {
		UserIDstr := chi.URLParam(r, "id")
		UserIDreal, errParse = uuid.Parse(UserIDstr)

	} else {
		UserIDreal, errParse = uuid.Parse(UserIDstri)
	}

	if errParse != nil {
		jsonresponse.RespondWithBadRequest(w, "gagl men parse ID dari context")
		return
	}

	app.Service.UpdateUser(r.Context(), UserIDreal, params.Nama, params.Password, params.Email, params.Role)
}

func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {

	UserIDstr := chi.URLParam(r, "id")

	UserID, erro := uuid.Parse(UserIDstr)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menparse %v", erro))
	}

	err := app.Service.DeleteUser(r.Context(), UserID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal meng hapus data %v", err))
	}
}

func (app *Application) GetListUser(w http.ResponseWriter, r *http.Request) {
	pagestr := r.URL.Query().Get("page")
	pagesizestr := r.URL.Query().Get("pagesize")

	Page := 1
	PageSize := 10

	if pagestr != "" {
		p, err := strconv.Atoi(pagestr)

		if err != nil && Page > 0 {
			Page = p
		}
	}

	if pagestr != "" {
		ps, err := strconv.Atoi(pagesizestr)

		if err != nil && ps > 0 {
			PageSize = ps
		}
	}

	User, erro := app.Service.ListsUserID(r.Context(), Page, PageSize)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan data %v", erro))
	}

	jsonresponse.ResponSuccess(w, 200, User)

}
