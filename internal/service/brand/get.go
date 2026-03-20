package brand

import (
	"context"
	"errors"
	"salon/internal/models"
	"time"

	repo "salon/internal/repository/errors"

	"github.com/apple5343/errorx"
	"github.com/docker/distribution/uuid"
)

const (
	ttl = 5 * time.Minute
)

func (s *brandService) getByID(ctx context.Context, id string) (*models.Brand, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	b, err := s.cache.GetByID(ctx, id)
	if nil == err {
		return b, nil
	}
	b, err = s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrBrandNotFound
		}
		return nil, errorx.NewError("get brand: "+err.Error(), errorx.Internal)
	}
	s.cache.SetByID(ctx, b, ttl) //TODO логи
	return b, nil
}

func (s *brandService) GetByID(ctx context.Context, id string) (*models.Brand, error) {
	return s.getByID(ctx, id)
}

func (s *brandService) GetBrands(ctx context.Context, filter *models.BrandFilters) ([]*models.Brand, error) {
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
	if filter.OrderBy == nil {
		orderBy := models.BrandOrderByCreatedAt
		filter.OrderBy = &orderBy
	}
	return s.repo.GetBrandsByFilter(ctx, filter)
}
