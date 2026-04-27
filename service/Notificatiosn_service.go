package service

import (
	"context"
	"database/sql"
	"errors"
	"pet-care/database"
	"pet-care/internal/middleware"

	"github.com/google/uuid"
)

func (S *Services) GetNotificicationsHistory(ctx context.Context, Page, PageSize int32) ([]database.GetHistoryNotificationsRow, error) {
	Role, okey := middleware.GetRoleFromContext(ctx)

	if !okey {
		return []database.GetHistoryNotificationsRow{}, errors.New("role tidak valid")
	}

	if Role == "Users" {
		return []database.GetHistoryNotificationsRow{}, errors.New("anda tidak memiliki akses")
	}

	UsersIDstr, oke := middleware.GetIDFromContext(ctx)

	if !oke {
		return []database.GetHistoryNotificationsRow{}, errors.New("gagal mendapatkan id di context")
	}

	UserID, errs := uuid.Parse(UsersIDstr)

	if errs != nil {
		return []database.GetHistoryNotificationsRow{}, errs
	}

	Offset := (Page - 1) * PageSize

	NotifHistory, err := S.StoreDB.Notifications.GetHistoryNotifications(ctx, database.GetHistoryNotificationsParams{
		IDUser: UserID,
		Offset: int32(Offset),
		Limit:  int32(PageSize),
	})

	if err != nil {
		return []database.GetHistoryNotificationsRow{}, err
	}

	return NotifHistory, err
}

func (S *Services) DeleteNotifications(ctx context.Context, ID uuid.UUID) error {
	role, okey := middleware.GetRoleFromContext(ctx)

	if !okey {
		return errors.New("role tidak valid")
	}

	if role == "User" {
		return errors.New("anda tidak memiliki hak akses ini")
	}

	err := S.StoreDB.Notifications.DeleteNotifications(ctx, ID)

	if err != nil {
		return err
	}
	return nil
}

func (S *Services) UpdateNotification(ctx context.Context, title, message string, id uuid.UUID) (database.Notification, error) {
	role, okey := middleware.GetRoleFromContext(ctx)

	if !okey {
		return database.Notification{}, errors.New("role tidak valid")
	}
	if role == "User" {
		return database.Notification{}, errors.New("anda tidak memiliki hak akses ini")
	}

	titles := sql.NullString{
		String: title,
		Valid:  true,
	}

	Notification, err := S.StoreDB.Notifications.UpdateNotification(ctx, database.UpdateNotificationParams{
		Title:   titles,
		Message: message,
		ID:      id,
	})

	if err != nil {
		return database.Notification{}, err
	}
	return Notification, nil
}
