package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Message interface {
	CreateMessage(ctx context.Context, arg database.CreateMessageParams) (database.Message, error)
	DeleteMessage(ctx context.Context, arg database.DeleteMessageParams) error
	GetHistoryPesan(ctx context.Context, arg database.GetHistoryPesanParams) ([]database.GetHistoryPesanRow, error)
	GetChatInbox(ctx context.Context, arg database.GetChatInboxParams) ([]database.GetChatInboxRow, error)
	MarkMessageAsRead(ctx context.Context, arg database.MarkMessageAsReadParams) error
	GetMessagebyIDuser(ctx context.Context, senderID uuid.UUID) (database.GetMessagebyIDuserRow, error)
}

type message struct {
	q *database.Queries
}

func (m *message) CreateMessage(ctx context.Context, arg database.CreateMessageParams) (database.Message, error) {
	return m.q.CreateMessage(ctx, arg)
}
func (m *message) DeleteMessage(ctx context.Context, arg database.DeleteMessageParams) error {
	return m.q.DeleteMessage(ctx, arg)
}
func (m *message) GetHistoryPesan(ctx context.Context, arg database.GetHistoryPesanParams) ([]database.GetHistoryPesanRow, error) {
	return m.q.GetHistoryPesan(ctx, arg)
}
func (m *message) GetChatInbox(ctx context.Context, arg database.GetChatInboxParams) ([]database.GetChatInboxRow, error) {
	return m.q.GetChatInbox(ctx, arg)
}
func (m *message) MarkMessageAsRead(ctx context.Context, arg database.MarkMessageAsReadParams) error {
	return m.q.MarkMessageAsRead(ctx, arg)
}

func (m *message) GetMessagebyIDuser(ctx context.Context, senderID uuid.UUID) (database.GetMessagebyIDuserRow, error) {
	return m.q.GetMessagebyIDuser(ctx, senderID)
}
