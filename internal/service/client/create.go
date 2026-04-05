package client

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *clientService) Register(ctx context.Context, c *models.Client) (*models.Client, error) {
	userRole := ctxutil.UserRoleFromContext(ctx)
	if userRole == "" {
		return nil, ErrUnauthorized
	} else if userRole != string(models.EmployeeRoleAdmin) && userRole != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	if err := c.BeforeCreate(s.clock); err != nil {
		return nil, errorx.Wrap("register client", errorx.BadRequest, err)
	}
	c, err := s.repo.Create(ctx, c)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, ErrClientExists
		}
		return nil, errorx.Wrap("create client", errorx.Internal, err)
	}
	//TODO добавить логи
	return c, nil
}
