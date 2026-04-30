package main

import (
	"errors"
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func HelperIDMessage(r *http.Request) (uuid.UUID, error) {
	MessageiDstr := chi.URLParam(r, "id_message")

	if MessageiDstr == "" {
		return uuid.Nil, errors.New("Not Fund")
	}

	MessageID, erro := uuid.Parse(MessageiDstr)

	if erro != nil {
		return uuid.Nil, erro
	}
	return MessageID, nil
}

func (app Application) GetChatInbox(w http.ResponseWriter, r *http.Request) {

	Useridstr, okey := middleware.GetIDFromContext(r.Context())

	if !okey {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID %v", erros))
		return
	}

	UserID, errs := uuid.Parse(Useridstr)

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal meng-parse ID %v", errs))
		return
	}
	page, pagezise, errr := HelperPage(r)

	if errr != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan info page %v", errr))
	}
	message, err := app.Service.GetChatInbox(r.Context(), UserID, page, pagezise)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan message %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, message)

}

func (app Application) GetChatHistory(w http.ResponseWriter, r *http.Request) {

	id_booking, errs := HelperIDBookings(r)

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan id_booking dari context %v", errs))
		return
	}

	page, pagesize, erro := HelperPage(r)

	if erro != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan info page %v", erro))
		return
	}
	chath, err := app.Service.GetChatHistory(r.Context(), id_booking, int32(page), int32(pagesize))

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan chat histore %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, 200, chath)
}

func (app Application) DeleteMessage(w http.ResponseWriter, r *http.Request) {

	IDMessage, errs := HelperIDMessage(r)

	if errs != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal mendapatkan ID message %v", errs))
	}
	err := app.Service.DeleteMessage(r.Context(), IDMessage)

	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal menghapus pesan %v", err))
		return
	}
}
