package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Pets interface {
	CreatePets(ctx context.Context, arg database.CreatePetsParams) (database.Pet, error)
	DeletePets(ctx context.Context, arg database.DeletePetsParams) error
	GetPetsListSt(ctx context.Context, arg database.GetPetsListStParams) ([]database.GetPetsListStRow, error)
	UpdatePets(ctx context.Context, arg database.UpdatePetsParams) (database.Pet, error)
	GetPetsListUser(ctx context.Context, arg database.GetPetsListUserParams) ([]database.GetPetsListUserRow, error)
	GetPetsDetail(ctx context.Context, id uuid.UUID) (database.GetPetsDetailRow, error)
	GetPetsByIDUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type pets struct {
	q *database.Queries
}

func (p *pets) CreatePets(ctx context.Context, arg database.CreatePetsParams) (database.Pet, error) {
	return p.q.CreatePets(ctx, arg)
}

func (p *pets) DeletePets(ctx context.Context, arg database.DeletePetsParams) error {
	return p.q.DeletePets(ctx, arg)
}

func (p *pets) GetPetsDetail(ctx context.Context, id uuid.UUID) (database.GetPetsDetailRow, error) {
	return p.q.GetPetsDetail(ctx, id)
}

func (p *pets) GetPetsListUser(ctx context.Context, arg database.GetPetsListUserParams) ([]database.GetPetsListUserRow, error) {
	return p.q.GetPetsListUser(ctx, arg)
}

func (p *pets) GetPetsByIDUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return p.q.GetPetsByIDUser(ctx, userID)
}

func (p *pets) GetPetsListSt(ctx context.Context, arg database.GetPetsListStParams) ([]database.GetPetsListStRow, error) {
	return p.q.GetPetsListSt(ctx, arg)
}

func (p *pets) UpdatePets(ctx context.Context, arg database.UpdatePetsParams) (database.Pet, error) {
	return p.q.UpdatePets(ctx, arg)
}
