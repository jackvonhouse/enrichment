package usecase

import (
	"github.com/jackvonhouse/enrichment/app/service"
	"github.com/jackvonhouse/enrichment/internal/usecase/user"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type UseCase struct {
	User user.UseCase
}

func New(
	service service.Service,
	logger log.Logger,
) UseCase {

	useCaseLogger := logger.WithField("layer", "usecase")

	return UseCase{
		User: user.New(
			service.Enrichment,
			service.User,
			useCaseLogger,
		),
	}
}
