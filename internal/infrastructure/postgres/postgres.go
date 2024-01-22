package postgres

import (
	"context"
	"fmt"
	"github.com/jackvonhouse/enrichment/config"
	"github.com/jackvonhouse/enrichment/pkg/log"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func New(
	ctx context.Context,
	config config.Database,
	logger log.Logger,
) (Database, error) {

	logger.Info(config.String())

	db, err := sqlx.ConnectContext(ctx, "postgres", config.String())

	if err != nil {
		logger.Warnf("can't connect to postgres: %s", err)

		return Database{}, fmt.Errorf("can't connect to postgres: %s", err)
	}

	return Database{
		db: db,
	}, nil
}

func (d Database) Database() *sqlx.DB { return d.db }
