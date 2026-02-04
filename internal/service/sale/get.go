package sale

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *saleService) getByID(ctx context.Context, id string) (*models.Sale, error) {
	sale, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrSaleNotFound
		}
		return nil, errorx.NewError("get sale: "+err.Error(), errorx.Internal)
	}
	return sale, nil
}

func (s *saleService) GetByID(ctx context.Context, id string) (*models.Sale, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	sale, err := s.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sale.EmployeeID != ctxutil.UserIDFromContext(ctx) && role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	return sale, nil
}
