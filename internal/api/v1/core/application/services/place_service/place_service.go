package place_service

import (
	"context"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/exceptions"
)

type PlaceService struct {
	pr IPlaceRepository
}

type IPlaceRepository interface {
	GetAll(ctx context.Context) (*[]dto.PlaceDto, error)
	Get(ctx context.Context, id int64) (*dto.PlaceDto, error)
	Create(ctx context.Context, p *dto.CreatePlaceDto) (*dto.PlaceDto, error)
	Delete(ctx context.Context, id int64) error
}

func NewPlaceService(pr IPlaceRepository) *PlaceService {
	return &PlaceService{pr: pr}
}

func (ps *PlaceService) GetPlaces(ctx context.Context) (*[]dto.PlaceDto, *exceptions.Error_) {
	places, err := ps.pr.GetAll(ctx)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	return places, nil
}

func (ps *PlaceService) PlaceIsExists(ctx context.Context, placeID int64) (*bool, *exceptions.Error_) {
	place, err := ps.pr.Get(ctx, placeID)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	state := place != nil
	return &state, nil
}

func (ps *PlaceService) CreatePlace(ctx context.Context, p *dto.CreatePlaceDto) (*dto.PlaceDto, *exceptions.Error_) {
	place, err := ps.pr.Create(ctx, p)
	if err != nil {
		return nil, &exceptions.ErrDatabaseError
	}
	return place, nil
}

func (ps *PlaceService) DeletePlace(ctx context.Context, id int64) *exceptions.Error_ {
	place, err := ps.pr.Get(ctx, id)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	if place == nil {
		return &exceptions.ErrPlaceNotFound
	}

	err = ps.pr.Delete(ctx, id)
	if err != nil {
		return &exceptions.ErrDatabaseError
	}
	return nil
}
