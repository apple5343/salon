package car

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"

	"github.com/apple5343/errorx"
)

func (s *carService) getCarByID(ctx context.Context, id string) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	//TODO кеширование
	c, err := s.repo.GetCarByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, nil, nil, nil, ErrCarNotFound
		}
		return nil, nil, nil, nil, errorx.NewError("get car: "+err.Error(), errorx.Internal)
	}
	m, b, err := s.model.GetByID(ctx, c.ModelID)
	if err != nil {
		return nil, nil, nil, nil, errorx.NewError("create car: "+err.Error(), errorx.Internal)
	}
	sup, err := s.supplier.GetByID(ctx, c.SupplierID)
	if err != nil {
		return nil, nil, nil, nil, errorx.NewError("create car: "+err.Error(), errorx.Internal)
	}
	return c, m, b, sup, nil
}

func (s *carService) GetCarByID(ctx context.Context, id string) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	return s.getCarByID(ctx, id)
}
