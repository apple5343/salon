package supplier

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	"salon/pkg/logger"
	"time"

	"github.com/apple5343/errorx"
	"github.com/google/uuid"
)

const (
	ttl = 5 * time.Minute
)

func (s *supplierService) getByID(ctx context.Context, id string) (*models.Supplier, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	supplier, err := s.cache.GetByID(ctx, id)
	if nil == err {
		return supplier, nil
	}
	supplier, err = s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrSupplierNotFound
		}
		return nil, errorx.NewError("get supplier: "+err.Error(), errorx.Internal)
	}
	if err = s.cache.SetByID(ctx, supplier, ttl); err != nil {
		logger.FromContextOrDefault(ctx).Error(ctx, "get supplier: "+err.Error())
	}
	return supplier, nil
}

func (s *supplierService) GetByID(ctx context.Context, id string) (*models.Supplier, error) {
	return s.getByID(ctx, id)
}
