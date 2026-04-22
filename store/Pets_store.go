package Store

import (
	"context"
	"pet-care/database"

	"github.com/google/uuid"
)

type Pets interface {
	CreatePets(ctx context.Context, arg database.CreatePetsParams) (database.Pet, error)
	DeletePets(ctx context.Context, arg database.DeletePetsParams) error
	GetPetsID(ctx context.Context, id uuid.UUID) (database.GetPetsIDRow, error)
	GetPetsMany(ctx context.Context, arg database.GetPetsManyParams) ([]database.GetPetsManyRow, error)
	UpdatePets(ctx context.Context, arg database.UpdatePetsParams) (database.Pet, error)
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

func (p *pets) GetPetsID(ctx context.Context, id uuid.UUID) (database.GetPetsIDRow, error) {
	return p.q.GetPetsID(ctx, id)
}

func (p *pets) GetPetsMany(ctx context.Context, arg database.GetPetsManyParams) ([]database.GetPetsManyRow, error) {
	return p.q.GetPetsMany(ctx, arg)
}

func (p *pets) UpdatePets(ctx context.Context, arg database.UpdatePetsParams) (database.Pet, error) {
	return p.q.UpdatePets(ctx, arg)
}
