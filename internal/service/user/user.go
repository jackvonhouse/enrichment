package user

import (
	"context"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type repositoryUser interface {
	Create(context.Context, dto.CreateDTO, dto.EnrichmentDTO) (int, error)

	Get(context.Context, dto.GetDTO) ([]dto.User, error)
	GetById(context.Context, int) (dto.User, error)
	GetByFilter(context.Context, dto.GetDTO, dto.FilterDTO) ([]dto.User, error)

	Update(context.Context, dto.UpdateDTO) (int, error)

	Delete(context.Context, int) (int, error)
}

type Service struct {
	repository repositoryUser
	logger     log.Logger
}

func New(
	logger log.Logger,
	repository repositoryUser,
) Service {

	return Service{
		repository: repository,
		logger:     logger.WithField("unit", "user"),
	}
}

func (s Service) Create(
	ctx context.Context,
	create dto.CreateDTO,
	enrichment dto.EnrichmentDTO,
) (int, error) {

	return s.repository.Create(ctx, create, enrichment)
}

func (s Service) Get(
	ctx context.Context,
	data dto.GetDTO,
) ([]dto.User, error) {

	return s.repository.Get(ctx, data)
}

func (s Service) GetById(
	ctx context.Context,
	id int,
) (dto.User, error) {

	return s.repository.GetById(ctx, id)
}

func (s Service) GetByFilter(
	ctx context.Context,
	get dto.GetDTO,
	filter dto.FilterDTO,
) ([]dto.User, error) {

	return s.repository.GetByFilter(ctx, get, filter)
}

func (s Service) Update(
	ctx context.Context,
	data dto.UpdateDTO,
) (int, error) {

	return s.repository.Update(ctx, data)
}

func (s Service) Delete(
	ctx context.Context,
	id int,
) (int, error) {

	return s.repository.Delete(ctx, id)
}
