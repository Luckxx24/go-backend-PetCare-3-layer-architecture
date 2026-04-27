package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/database"
	"pet-care/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func Helperrolenotif(r *http.Request) (uuid.UUID, error) {
	role, oke := middleware.GetRoleFromContext(r.Context())

	if !oke {
		return uuid.Nil, errors.New("gagal mendapatkan role dari middleware")
	}

	if role != "Admin " || role != "Staff" {
		return uuid.Nil, errors.New("Unathorized")
	}

	UserIDstr := chi.URLParam(r, "id_user")

	if UserIDstr == "" {
		return uuid.Nil, errors.New("gagal menemukan id di url param")
	}

	UserID, err := uuid.Parse(UserIDstr)

	if err != nil {
		return uuid.Nil, err
	}

	return UserID, nil
}

type paramss struct {
	Title   string
	Message string
}

func (app Application) CreateNotifications(w http.ResponseWriter, r *http.Request) {
	type paramss struct {
	}

	app.Service.StoreDB.Notifications.CreateNotifications(r.Context(), database.CreateNotificationsParams{})

}

func (app Application) DeleteNofications(w http.ResponseWriter, r *http.Request) {

	UserID, erro := Helperrole(r)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mengambil id user dari url %v", erro))
		return
	}
	err := app.Service.DeleteNotifications(r.Context(), UserID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menghapus data %v", err))
		return
	}

}

func (app Application) UpdateNotification(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	param := paramss{}
	decode.Decode(&param)

	NOtifIDstr := chi.URLParam(r, "id_notifications")

	if NOtifIDstr == "" {
		jsonresponse.RespondWithNotfound(w, "gagal menemukan id notifiction dari url")
		return
	}

	NotifID, errr := uuid.Parse(NOtifIDstr)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal men-parse ID %v", errr))
		return
	}

	UpdatedNotifications, err := app.Service.UpdateNotification(r.Context(), param.Title, param.Message, NotifID)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal meng update Notifications %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, UpdatedNotifications)
}

func (app Application) GetNotification(w http.ResponseWriter, r *http.Request) {
	app.Service.Ge
}
