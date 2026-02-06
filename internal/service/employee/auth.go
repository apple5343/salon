package employee

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"
	utils "salon/internal/utils/password"
	"salon/pkg/jwt"

	"github.com/apple5343/errorx"
)

func (s *employeeService) Login(ctx context.Context, email, password string) (string, string, error) {
	e, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return "", "", ErrInvalidCreds
		}
		return "", "", errorx.NewError("login employee: "+err.Error(), errorx.Internal)
	}
	if e.Status != models.EmployeeStatusActive {
		return "", "", ErrInvalidCreds
	}
	if !utils.CheckPasswordHash(password, e.PasswordHash) {
		return "", "", ErrInvalidCreds
	}

	refreshToken, err := jwt.GenerateToken(jwt.UserInfo{ID: e.ID, Role: string(e.Role)}, []byte(s.jwtConfig.RefreshSecret), s.jwtConfig.RefreshTTL)
	if err != nil {
		return "", "", errorx.NewError("login employee: "+err.Error(), errorx.Internal)
	}

	accessToken, err := jwt.GenerateToken(jwt.UserInfo{ID: e.ID, Role: string(e.Role)}, []byte(s.jwtConfig.AccessSecret), s.jwtConfig.AccessTTL)
	if err != nil {
		return "", "", errorx.NewError("login employee: "+err.Error(), errorx.Internal)
	}

	return accessToken, refreshToken, nil
}

func (s *employeeService) GetRefreshToken(ctx context.Context) (string, error) {
	refreshToken := ctxutil.UserTokenFromContext(ctx)
	if refreshToken == "" {
		return "", ErrInvalidToken
	}
	userCalims, err := jwt.VerifyToken(refreshToken, []byte(s.jwtConfig.RefreshSecret))
	if err != nil {
		return "", ErrInvalidToken
	}
	e, err := s.getByID(ctx, userCalims.ID)
	if err != nil {
		if errors.Is(err, ErrEmployeeNotFound) {
			return "", ErrInvalidToken
		}
		return "", err
	}
	if e.Status != models.EmployeeStatusActive {
		return "", ErrInvalidToken
	}
	refreshToken, err = jwt.GenerateToken(jwt.UserInfo{ID: e.ID, Role: string(e.Role)}, []byte(s.jwtConfig.RefreshSecret), s.jwtConfig.RefreshTTL)
	if err != nil {
		return "", errorx.NewError("get refreshToken employee: "+err.Error(), errorx.Internal)
	}

	return refreshToken, nil
}

func (s *employeeService) GetAccessToken(ctx context.Context) (string, error) {
	refreshToken := ctxutil.UserTokenFromContext(ctx)
	if refreshToken == "" {
		return "", ErrInvalidToken
	}
	userCalims, err := jwt.VerifyToken(refreshToken, []byte(s.jwtConfig.RefreshSecret))
	if err != nil {
		return "", ErrInvalidToken
	}
	e, err := s.getByID(ctx, userCalims.ID)
	if err != nil {
		if errors.Is(err, ErrEmployeeNotFound) {
			return "", ErrInvalidToken
		}
		return "", err
	}
	if e.Status != models.EmployeeStatusActive {
		return "", ErrInvalidToken
	}
	accessToken, err := jwt.GenerateToken(jwt.UserInfo{ID: e.ID, Role: string(e.Role)}, []byte(s.jwtConfig.AccessSecret), s.jwtConfig.AccessTTL)
	if err != nil {
		return "", errorx.NewError("get accessToken employee: "+err.Error(), errorx.Internal)
	}

	return accessToken, nil
}
