package service

import (
	"context"
	"database/sql"
	"errors"
	"pet-care/database"
	"pet-care/internal/middleware"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (S *Services) CreateMessage(ctx context.Context, message string, reveiver_id, bookingsID uuid.UUID) (database.Message, error) {
	User_IDstr, ok := middleware.GetIDFromContext(ctx)

	if !ok {
		return database.Message{}, errors.New("gagal mengambil id dari context")
	}

	User_ID, erro := uuid.Parse(User_IDstr)
	if erro != nil {
		return database.Message{}, erro
	}

	bookings, err := S.StoreDB.Bookings.GetBookingByUserID(ctx, database.GetBookingByUserIDParams{
		ID:     bookingsID,
		UserID: User_ID,
	})

	if err != nil {
		return database.Message{}, err
	}

	message = strings.TrimSpace(message)

	if message == "" {
		return database.Message{}, errors.New("Harap masukan pesan")
	}

	waktu := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	Message, erro := S.StoreDB.Message.CreateMessage(ctx, database.CreateMessageParams{
		ID:         uuid.New(),
		BookingsID: bookings.ID,
		SenderID:   User_ID,
		ReceiverID: reveiver_id,
		Message:    message,
		IsRead:     false,
		CreatedAt:  waktu,
	})

	if erro != nil {
		return database.Message{}, erro
	}

	return Message, nil
}

func (S *Services) GetChatInbox(ctx context.Context, UserID uuid.UUID, Page, PageSize int) ([]database.GetChatInboxRow, error) {

	User_ID, ok := middleware.GetIDFromContext(ctx)

	if !ok {
		return []database.GetChatInboxRow{}, errors.New("gagal mengambil id dari context")
	}

	SenderID, err := uuid.Parse(User_ID)

	if err != nil {
		return []database.GetChatInboxRow{}, err
	}

	Offset := (Page - 1) * PageSize

	GetChatin, err := S.StoreDB.Message.GetChatInbox(ctx, database.GetChatInboxParams{
		SenderID: SenderID,
		Offset:   int32(Offset),
		Limit:    int32(PageSize),
	})

	if err != nil {
		return []database.GetChatInboxRow{}, err
	}

	return GetChatin, nil

}

func (S *Services) GetChatHistory(ctx context.Context, id_booking uuid.UUID, Page, PageSize int32) ([]database.GetHistoryPesanRow, error) {
	UserIDstr, ok := middleware.GetIDFromContext(ctx)

	if !ok {
		return []database.GetHistoryPesanRow{}, errors.New("gagal mengambil id dari context")
	}

	UserID, erro := uuid.Parse(UserIDstr)

	if erro != nil {
		return []database.GetHistoryPesanRow{}, erro
	}

	Bookings, errs := S.StoreDB.Bookings.GetBookingByUserID(ctx, database.GetBookingByUserIDParams{
		ID:     id_booking,
		UserID: UserID,
	})

	if errs != nil {
		return []database.GetHistoryPesanRow{}, erro
	}

	Role, okey := middleware.GetRoleFromContext(ctx)

	if !okey {
		return []database.GetHistoryPesanRow{}, errors.New("gagal mengambil role dari context")
	}

	oke := IsValidRole(Role)

	if !oke {
		return []database.GetHistoryPesanRow{}, errors.New("role tidak valid")
	}

	if Role == "User" {
		return []database.GetHistoryPesanRow{}, errors.New("anda tidak memiliki hak ini")
	}

	Offset := (Page - 1) * PageSize

	ChatHistory, err := S.StoreDB.Message.GetHistoryPesan(ctx, database.GetHistoryPesanParams{
		BookingsID: Bookings.ID,
		Offset:     Offset,
		Limit:      PageSize,
	})

	if err != nil {
		return []database.GetHistoryPesanRow{}, err
	}

	return ChatHistory, nil
}

func (S *Services) DeleteMessage(ctx context.Context, ID uuid.UUID) (bool, error) {

	id_USERstr, ok := middleware.GetIDFromContext(ctx)

	if !ok {
		return false, errors.New("gagal mengambil id dari context")
	}
	Users_ID, err := uuid.Parse(id_USERstr)

	if err != nil {
		return false, err
	}

	erro := S.StoreDB.Message.DeleteMessage(ctx, database.DeleteMessageParams{
		ID:       ID,
		SenderID: Users_ID,
	})

	if erro != nil {
		return false, erro
	}

	return true, nil

}

func (S *Services) UpdateMessageAsRead(ctx context.Context, booking_id uuid.UUID) error {
	UserIDstr, okey := middleware.GetIDFromContext(ctx)

	if !okey {
		return errors.New("gagal mengambil id dari context")
	}

	User_ID, _ := uuid.Parse(UserIDstr)

	booking, err := S.StoreDB.Bookings.GetBookingByUserID(ctx, database.GetBookingByUserIDParams{
		ID:     booking_id,
		UserID: User_ID,
	})

	if err != nil {
		return err
	}

	Error := S.StoreDB.Message.MarkMessageAsRead(ctx, database.MarkMessageAsReadParams{
		BookingsID: booking.ID,
		SenderID:   User_ID,
	})

	if Error != nil {
		return Error
	}

	return nil
}
