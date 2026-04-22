package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Notifications interface {
	CreateNotifications(ctx context.Context, arg database.CreateNotificationsParams) (database.Notification, error)
	DeleteNotifications(ctx context.Context, id uuid.UUID) error
	GetHistoryNotifications(ctx context.Context, arg database.GetHistoryNotificationsParams) ([]database.GetHistoryNotificationsRow, error)
	GetNotifications(ctx context.Context, id uuid.UUID) (database.GetNotificationsRow, error)
	UpdateNotification(ctx context.Context, arg database.UpdateNotificationParams) (database.Notification, error)
}

type notifications struct {
	q *database.Queries
}

func (n notifications) CreateNotifications(ctx context.Context, arg database.CreateNotificationsParams) (database.Notification, error) {
	return n.q.CreateNotifications(ctx, arg)
}

func (n notifications) DeleteNotifications(ctx context.Context, id uuid.UUID) error {
	return n.q.DeleteNotifications(ctx, id)
}

func (n notifications) GetHistoryNotifications(ctx context.Context, arg database.GetHistoryNotificationsParams) ([]database.GetHistoryNotificationsRow, error) {
	return n.q.GetHistoryNotifications(ctx, arg)
}

func (n notifications) GetNotifications(ctx context.Context, id uuid.UUID) (database.GetNotificationsRow, error) {
	return n.q.GetNotifications(ctx, id)
}

func (n notifications) UpdateNotification(ctx context.Context, arg database.UpdateNotificationParams) (database.Notification, error) {
	return n.q.UpdateNotification(ctx, arg)
}
