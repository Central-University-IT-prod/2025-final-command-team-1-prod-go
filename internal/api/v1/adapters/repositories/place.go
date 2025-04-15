package repositories

import (
	"context"
	"database/sql"
	"errors"

	"example.com/m/internal/api/v1/core/application/dto"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type PlaceRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPlaceRepository(db *sql.DB, logger *zap.Logger) *PlaceRepository {
	return &PlaceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PlaceRepository) GetAll(ctx context.Context) (*[]dto.PlaceDto, error) {
	var places []dto.PlaceDto
	query, _, _ := goqu.From("places").ToSQL()

	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return &places, nil
	}
	if err != nil {
		r.logger.Error(
			"Place Repository Error",
			zap.String("method", "GetAll"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	for rows.Next() {
		var place dto.PlaceDto

		rows.Scan(&place.ID, &place.Name, &place.Description, &place.Address, &place.City)
		places = append(places, place)
	}

	return &places, nil
}

func (r *PlaceRepository) Get(ctx context.Context, id int64) (*dto.PlaceDto, error) {
	var place dto.PlaceDto
	query, _, _ := goqu.From("places").Where(goqu.Ex{
		"id": id,
	}).ToSQL()

	err := r.db.QueryRow(query).Scan(&place.ID, &place.Name, &place.Description, &place.Address, &place.City)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.logger.Error(
			"Place Repository Error",
			zap.String("method", "Get"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &place, nil
}

func (r *PlaceRepository) Delete(ctx context.Context, id int64) error {
	query, _, _ := goqu.From("places").Where(goqu.C("id").Eq(id)).Delete().ToSQL()
	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Place Repository Error",
			zap.String("method", "Delete"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *PlaceRepository) Create(ctx context.Context, p *dto.CreatePlaceDto) (*dto.PlaceDto, error) {
	query, _, err := goqu.
		Insert("places").
		Rows(*p).
		Returning("id").
		ToSQL()

	if err != nil {
		r.logger.Error(
			"Place Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	var id int64
	err = r.db.QueryRow(query).Scan(&id)
	if err != nil {
		r.logger.Error(
			"Place Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	x := p.ToPlaceDto(id)
	return &x, nil
}
