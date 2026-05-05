package main

import (
	"fmt"
	"net/http"
	"pet-care/cmd/jsonresponse"
	"pet-care/cmd/ws"
	"pet-care/internal/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app Application) CreateMessage(w http.ResponseWriter, r *http.Request) {

	bookingID, err := HelperIDBookings(r)
	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("booking_id tidak valid: %v", err))
		return
	}

	var body struct {
		ReceiverID uuid.UUID `json:"receiver_id"`
		Message    string    `json:"message"`
	}

	if err := readJSON(r, &body); err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("request body tidak valid: %v", err))
		return
	}

	msg, err := app.Service.CreateMessage(r.Context(), body.Message, body.ReceiverID, bookingID)
	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal membuat pesan: %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, http.StatusCreated, msg)
}

func (app Application) UpdateMessageAsRead(w http.ResponseWriter, r *http.Request) {
	bookingID, err := HelperIDBookings(r)
	if err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("booking_id tidak valid: %v", err))
		return
	}

	if err := app.Service.UpdateMessageAsRead(r.Context(), bookingID); err != nil {
		jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("gagal update status baca: %v", err))
		return
	}

	jsonresponse.ResponSuccess(w, http.StatusOK, map[string]string{
		"message": "pesan berhasil ditandai sebagai dibaca",
	})
}

func (app Application) ServeWS(hub *ws.MessageHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIDstr, ok := middleware.GetIDFromContext(r.Context())
		if !ok {
			jsonresponse.RespondWithBadRequest(w, "gagal mendapatkan user ID dari token")
			return
		}

		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			jsonresponse.RespondWithBadRequest(w, "user ID tidak valid")
			return
		}

		bookingID, err := HelperIDBookings(r)
		if err != nil {
			jsonresponse.RespondWithBadRequest(w, fmt.Sprintf("booking_id tidak valid: %v", err))
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {

			fmt.Printf("[WS] Gagal upgrade koneksi untuk user %v: %v\n", userID, err)
			return
		}

		client := &ws.UserConnection{
			UserID:    userID,
			BookingID: bookingID,
			Conn:      conn,

			Send:    make(chan ws.WsResponse, 256),
			Hub:     hub,
			Service: &app.Service,
		}

		hub.UserJoined <- client

		go client.TerimaPesan()

		client.KirimPesan()
	}
}
