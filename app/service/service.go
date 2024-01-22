package service

import (
	"github.com/jackvonhouse/enrichment/app/repository"
	"github.com/jackvonhouse/enrichment/internal/service/enrichment"
	"github.com/jackvonhouse/enrichment/internal/service/user"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type Service struct {
	Enrichment enrichment.Service
	User       user.Service
}

func New(
	repository repository.Repository,
	logger log.Logger,
) Service {

	serviceLogger := logger.WithField("layer", "service")

	return Service{
		Enrichment: enrichment.New(serviceLogger),
		User:       user.New(serviceLogger, repository.User),
	}
}
