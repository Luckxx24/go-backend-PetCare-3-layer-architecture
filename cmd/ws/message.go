package ws

import "github.com/google/uuid"

type WsPayload struct {
	Action     string    `json:"action"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Message    string    `json:"message"`
}

type WsResponse struct {
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Message    string    `json:"message"`
	BookingID  uuid.UUID `json:"booking_id"`
	Error      string    `json:"error,omitempty"`
}
