package analytics

import (
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden = errorx.NewError("forbidden", errorx.Forbidden)
)

type analyticsService struct {
	analyticsRepository repository.AnalyticsRepository
}

func NewService(repo repository.AnalyticsRepository) service.AnalyticsService {
	return &analyticsService{
		analyticsRepository: repo,
	}
}
