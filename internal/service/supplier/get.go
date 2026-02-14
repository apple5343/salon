package supplier

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"

	"github.com/apple5343/errorx"
	"github.com/docker/distribution/uuid"
)

func (s *supplierService) GetByID(ctx context.Context, id string) (*models.Supplier, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrSupplierNotFound
		}
		return nil, errorx.NewError("get supplier: "+err.Error(), errorx.Internal)
	}
	return supplier, nil
}
