package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Pet_Status_Log interface {
	CreateNewLog(ctx context.Context, arg database.CreateNewLogParams) (database.PetStatusLog, error)
	DeleteLog(ctx context.Context, id uuid.UUID) error
	GetAllLog(ctx context.Context, arg database.GetAllLogParams) ([]database.GetAllLogRow, error)
	GetLog(ctx context.Context, id uuid.UUID) (database.GetLogRow, error)
	UpdateLog(ctx context.Context, arg database.UpdateLogParams) (database.PetStatusLog, error)
	GetLOgbyIDbooking(ctx context.Context, idBookings uuid.UUID) (database.GetLOgbyIDbookingRow, error)
}

type pet_status_log struct {
	q *database.Queries
}

func (psl pet_status_log) CreateNewLog(ctx context.Context, arg database.CreateNewLogParams) (database.PetStatusLog, error) {
	return psl.q.CreateNewLog(ctx, arg)
}

func (psl pet_status_log) DeleteLog(ctx context.Context, id uuid.UUID) error {
	return psl.q.DeleteLog(ctx, id)
}

func (psl pet_status_log) GetAllLog(ctx context.Context, arg database.GetAllLogParams) ([]database.GetAllLogRow, error) {
	return psl.q.GetAllLog(ctx, arg)
}

func (psl pet_status_log) GetLog(ctx context.Context, id uuid.UUID) (database.GetLogRow, error) {
	return psl.q.GetLog(ctx, id)
}

func (psl pet_status_log) UpdateLog(ctx context.Context, arg database.UpdateLogParams) (database.PetStatusLog, error) {
	return psl.q.UpdateLog(ctx, arg)
}

func (psl pet_status_log) GetLOgbyIDbooking(ctx context.Context, idBookings uuid.UUID) (database.GetLOgbyIDbookingRow, error) {
	return psl.GetLOgbyIDbooking(ctx, idBookings)
}
