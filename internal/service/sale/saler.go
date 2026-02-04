package sale

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
	"github.com/shopspring/decimal"
)

func (s *saleService) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	car, _, _, _, err := s.carService.GetCarByID(ctx, sale.CarID)
	if err != nil {
		return nil, err
	}
	if car.Status != models.CarStatusAvailable {
		return nil, errorx.NewError("car is not available", errorx.BadRequest)
	}
	employeeID := ctxutil.UserIDFromContext(ctx)
	sale.EmployeeID = employeeID
	sale.Status = models.SaleStatusPending
	if err := sale.BeforeCreate(); err != nil {
		return nil, errorx.NewError(err.Error(), errorx.BadRequest)
	}
	sale.OriginPrice = car.Price
	if sale.DiscountAmount.GreaterThan(car.Price) {
		return nil, errorx.NewError("discount amount is greater than car price", errorx.BadRequest)
	}
	sale.FinalPrice = car.Price.Sub(sale.DiscountAmount)
	sale.DiscountPercent = sale.DiscountAmount.Div(car.Price).Mul(decimal.NewFromInt(100)).Round(2)
	return s.repo.Create(ctx, sale)
}

func (s *saleService) Complete(ctx context.Context, id string) (*models.Sale, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	employeeID := ctxutil.UserIDFromContext(ctx)
	sale, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sale.EmployeeID != employeeID && role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if sale.Status != models.SaleStatusPending {
		return nil, errorx.NewError("sale is not pending", errorx.BadRequest)
	}
	sale.Status = models.SaleStatusCompleted
	if err := sale.BeforeUpdate(); err != nil {
		return nil, err
	}
	err = s.repo.Complete(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, errorx.NewError("sale already processed or not found", errorx.BadRequest)
		}
		return nil, errorx.NewError("complete sale: "+err.Error(), errorx.Internal)
	}
	return s.getByID(ctx, id)
}

func (s *saleService) Cancel(ctx context.Context, id string) error {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return ErrForbidden
	}
	employeeID := ctxutil.UserIDFromContext(ctx)
	sale, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if sale.EmployeeID != employeeID && role != string(models.EmployeeRoleAdmin) {
		return ErrForbidden
	}
	sale.Status = models.SaleStatusCanceled
	if err := sale.BeforeUpdate(); err != nil {
		return err
	}
	if err := s.repo.Cancel(ctx, id); err != nil {
		return errorx.NewError("cancel sale: "+err.Error(), errorx.Internal)
	}
	return nil
}
