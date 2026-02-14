package client

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
	"github.com/docker/distribution/uuid"
)

func (s *clientService) getByID(ctx context.Context, id string) (*models.Client, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, errorx.NewError("get client: "+err.Error(), errorx.Internal)
	}
	return client, nil
}

func (s *clientService) GetByID(ctx context.Context, id string) (*models.Client, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	client, err := s.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return client, nil
}
