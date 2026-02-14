package car

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"

	"github.com/apple5343/errorx"
	"github.com/google/uuid"
)

func (s *carService) getBrandByID(ctx context.Context, id string) (*models.Brand, error) {
	//TODO кеширование
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	b, err := s.repo.GetBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrBrandNotFound
		}
		return nil, errorx.NewError("get brand: "+err.Error(), errorx.Internal)
	}
	return b, nil
}

func (s *carService) GetBrandByID(ctx context.Context, id string) (*models.Brand, error) {
	return s.getBrandByID(ctx, id)
}

func (s *carService) getModelByID(ctx context.Context, id string) (*models.Model, error) {
	//TODO кеширование
	m, err := s.repo.GetModelByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, errorx.NewError("get model: "+err.Error(), errorx.Internal)
	}
	return m, nil
}

func (s *carService) GetModelByID(ctx context.Context, id string) (*models.Model, *models.Brand, error) {
	model, err := s.getModelByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	brand, err := s.getBrandByID(ctx, model.BrandID)
	if err != nil {
		return nil, nil, err
	}
	return model, brand, nil
}

func (s *carService) getCarByID(ctx context.Context, id string) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	//TODO кеширование
	c, err := s.repo.GetCarByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, nil, nil, nil, ErrCarNotFound
		}
		return nil, nil, nil, nil, errorx.NewError("get car: "+err.Error(), errorx.Internal)
	}
	m, err := s.getModelByID(ctx, c.ModelID)
	if err != nil {
		return nil, nil, nil, nil, errorx.NewError("create car: "+err.Error(), errorx.Internal)
	}
	b, err := s.getBrandByID(ctx, m.BrandID)
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
