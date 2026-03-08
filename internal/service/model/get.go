package model

import (
	"context"
	"errors"
	"salon/internal/models"

	repo "salon/internal/repository/errors"

	"github.com/apple5343/errorx"
	"github.com/docker/distribution/uuid"
)

func (s *modelService) getModelByID(ctx context.Context, id string) (*models.Model, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	//TODO кеширование
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, errorx.NewError("get model: "+err.Error(), errorx.Internal)
	}
	return m, nil
}

func (s *modelService) GetByID(ctx context.Context, id string) (*models.Model, *models.Brand, error) {
	model, err := s.getModelByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	brand, err := s.brandService.GetByID(ctx, model.BrandID)
	if err != nil {
		return nil, nil, err
	}
	return model, brand, nil
}

func (s *modelService) GetModels(ctx context.Context, filter *models.ModelFilters) ([]*models.ModelShort, error) {
	if err := filter.Validate(); err != nil {
		return nil, errorx.NewError(err.Error(), errorx.BadRequest)
	}
	if filter.Limit == nil {
		limit := 10
		filter.Limit = &limit
	}
	if filter.Offset == nil {
		offset := 0
		filter.Offset = &offset
	}
	return s.repo.GetModelsByFilter(ctx, filter)
}
