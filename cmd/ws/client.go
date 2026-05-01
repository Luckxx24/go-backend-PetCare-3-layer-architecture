package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type Client struct {
	UserID    uuid.UUID
	BookingID uuid.UUID
	Conn      *websocket.Conn
	Send      chan WsResponse
	Hub       *Hub
	Service   MessageServiceIface
}

type MessageServiceIface interface {
	CreateMessageWS(userID, receiverID, bookingID uuid.UUID, message string) error
}

func (c *Client) ReadPump() {

	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[CLIENT] Error membaca pesan dari %v: %v", c.UserID, err)
			}
			break
		}

		var payload WsPayload
		if err := json.Unmarshal(rawMsg, &payload); err != nil {
			log.Printf("[CLIENT] Gagal decode payload dari %v: %v", c.UserID, err)

			c.Send <- WsResponse{Error: "format pesan tidak valid"}
			continue
		}

		switch payload.Action {

		case "send_message":

			if payload.Message == "" {
				c.Send <- WsResponse{Error: "pesan tidak boleh kosong"}
				continue
			}

			if err := c.Service.CreateMessageWS(c.UserID, payload.ReceiverID, c.BookingID, payload.Message); err != nil {
				log.Printf("[CLIENT] Gagal simpan pesan: %v", err)
				c.Send <- WsResponse{Error: "gagal menyimpan pesan"}
				continue
			}

			c.Hub.Broadcast <- BroadcastMsg{
				RoomID: c.BookingID,
				Message: WsResponse{
					SenderID:   c.UserID,
					ReceiverID: payload.ReceiverID,
					Message:    payload.Message,
					BookingID:  c.BookingID,
				},
			}

		default:
			log.Printf("[CLIENT] Action tidak dikenal dari %v: %s", c.UserID, payload.Action)
		}
	}
}

func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {

				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			jsonBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("[CLIENT] Gagal encode response: %v", err)
				return
			}

			w.Write(jsonBytes)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				nextMsg := <-c.Send
				nextBytes, _ := json.Marshal(nextMsg)
				w.Write(nextBytes)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
