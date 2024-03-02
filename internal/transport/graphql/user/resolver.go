package user

import (
	"context"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type useCaseUser interface {
	Create(context.Context, dto.CreateDTO) (int, error)

	Get(context.Context, dto.GetDTO, dto.FilterDTO, dto.SortDTO) ([]dto.User, error)
	GetById(context.Context, int) (dto.User, error)

	Update(context.Context, dto.UpdateDTO) (int, error)

	Delete(context.Context, int) (int, error)
}

type Resolver struct {
	useCase useCaseUser

	logger log.Logger
}

func New(
	useCase useCaseUser,
	logger log.Logger,
) Resolver {

	return Resolver{
		useCase: useCase,
		logger:  logger.WithField("transport_type", "graphql"),
	}
}
