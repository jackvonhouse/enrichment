package repository

import (
	"context"
	"github.com/jackvonhouse/enrichment/app/infrastructure"
	"github.com/jackvonhouse/enrichment/internal/infrastructure/postgres"
	"github.com/jackvonhouse/enrichment/internal/repository/user"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type Repository struct {
	User user.Repository

	Storage postgres.Database
}

func New(
	infrastructure infrastructure.Infrastructure,
	logger log.Logger,
) Repository {

	repositoryLogger := logger.WithField("layer", "repository")

	return Repository{
		User: user.New(
			infrastructure.Storage.Database(),
			repositoryLogger,
		),
		Storage: infrastructure.Storage,
	}
}

func (r Repository) Shutdown(
	_ context.Context,
) error {

	return r.Storage.Database().Close()
}
