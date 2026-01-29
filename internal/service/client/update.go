package client

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *clientService) Update(ctx context.Context, c *models.Client) (*models.Client, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	if err := c.BeforeUpdate(); err != nil {
		return nil, errorx.NewError("update client: "+err.Error(), errorx.BadRequest)
	}
	client, err := s.repo.Update(ctx, c)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrClientNotFound
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, ErrClientExists
		}
		return nil, errorx.NewError("update client: "+err.Error(), errorx.Internal)
	}
	return client, nil
}
