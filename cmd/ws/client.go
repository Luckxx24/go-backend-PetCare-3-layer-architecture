package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	batasWaktuMenulis = 10 * time.Second

	batasWaktuPong = 60 * time.Second

	intervalPing = (batasWaktuPong * 9) / 10

	ukuranMaksimalPesan = 512
)

type UserConnection struct {
	UserID    uuid.UUID
	BookingID uuid.UUID
	Conn      *websocket.Conn
	Send      chan WsResponse
	Hub       *MessageHub
	Service   MessageServiceIface
}

type MessageServiceIface interface {
	CreateMessageWS(userID, receiverID, bookingID uuid.UUID, message string) error
}

func (u *UserConnection) TerimaPesan() {
	defer func() {
		u.Hub.UserLeft <- u
		u.Conn.Close()
	}()

	u.Conn.SetReadLimit(ukuranMaksimalPesan)
	u.Conn.SetReadDeadline(time.Now().Add(batasWaktuPong))

	u.Conn.SetPongHandler(func(string) error {
		u.Conn.SetReadDeadline(time.Now().Add(batasWaktuPong))
		return nil
	})

	for {
		_, pesanMentah, err := u.Conn.ReadMessage()
		if err != nil {
			koneksiDitutupTidakNormal := websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			)
			if koneksiDitutupTidakNormal {
				log.Printf("[KONEKSI] Koneksi pengguna %v terputus tidak normal: %v", u.UserID, err)
			}
			break
		}

		var payload WsPayload
		if err := json.Unmarshal(pesanMentah, &payload); err != nil {
			log.Printf("[KONEKSI] Format pesan dari pengguna %v tidak valid: %v", u.UserID, err)
			u.Send <- WsResponse{Error: "format pesan tidak valid"}
			continue
		}

		switch payload.Action {

		case "send_message":
			if payload.Message == "" {
				u.Send <- WsResponse{Error: "pesan tidak boleh kosong"}
				continue
			}

			err := u.Service.CreateMessageWS(u.UserID, payload.ReceiverID, u.BookingID, payload.Message)
			if err != nil {
				log.Printf("[KONEKSI] Gagal menyimpan pesan dari pengguna %v: %v", u.UserID, err)
				u.Send <- WsResponse{Error: "gagal menyimpan pesan"}
				continue
			}

			u.Hub.NewMessage <- IncomingBroadcast{
				TargetRoomID: u.BookingID,
				Payload: WsResponse{
					SenderID:   u.UserID,
					ReceiverID: payload.ReceiverID,
					Message:    payload.Message,
					BookingID:  u.BookingID,
				},
			}

		default:
			log.Printf("[KONEKSI] Aksi tidak dikenal dari pengguna %v: %s", u.UserID, payload.Action)
		}
	}
}

func (u *UserConnection) KirimPesan() {
	timerPing := time.NewTicker(intervalPing)

	defer func() {
		timerPing.Stop()
		u.Conn.Close()
	}()

	for {
		select {

		case pesan, saluranMasihTerbuka := <-u.Send:
			u.Conn.SetWriteDeadline(time.Now().Add(batasWaktuMenulis))

			if !saluranMasihTerbuka {
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			penulis, err := u.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			pesanJSON, err := json.Marshal(pesan)
			if err != nil {
				log.Printf("[KONEKSI] Gagal encode pesan untuk pengguna %v: %v", u.UserID, err)
				return
			}
			penulis.Write(pesanJSON)

			jumlahPesanAntrian := len(u.Send)
			for i := 0; i < jumlahPesanAntrian; i++ {
				pesanBerikutnya := <-u.Send
				pesanBerikutnyaJSON, _ := json.Marshal(pesanBerikutnya)
				penulis.Write([]byte("\n"))
				penulis.Write(pesanBerikutnyaJSON)
			}

			if err := penulis.Close(); err != nil {
				return
			}

		case <-timerPing.C:
			u.Conn.SetWriteDeadline(time.Now().Add(batasWaktuMenulis))
			if err := u.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
