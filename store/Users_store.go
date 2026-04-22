package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Users interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) //all
	DeleteUser(ctx context.Context, id uuid.UUID) error                                   //admin
	GetUserID(ctx context.Context, id uuid.UUID) (database.GetUserIDRow, error)           //all
	GetUseremail(ctx context.Context, email string) (database.GetUseremailRow, error)     //all
	ListsUserID(ctx context.Context, arg database.ListsUserParams) ([]database.ListsUserRow, error)
	UpdateUser(ctx context.Context, arg database.UpdateUserParams) (database.User, error) //admin
}

type users struct {
	q *database.Queries
}

func (u *users) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return u.q.CreateUser(ctx, arg)
}

func (u *users) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.q.DeleteUser(ctx, id)
}

func (u *users) GetUserID(ctx context.Context, id uuid.UUID) (database.GetUserIDRow, error) {
	return u.q.GetUserID(ctx, id)
}

func (u *users) GetUseremail(ctx context.Context, email string) (database.GetUseremailRow, error) {
	return u.q.GetUseremail(ctx, email)
}

func (u *users) ListsUserID(ctx context.Context, arg database.ListsUserParams) ([]database.ListsUserRow, error) {
	return u.q.ListsUser(ctx, arg)
}

func (u *users) UpdateUser(ctx context.Context, arg database.UpdateUserParams) (database.User, error) {
	return u.q.UpdateUser(ctx, arg)
}
