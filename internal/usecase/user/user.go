package user

import (
	"context"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type serviceUser interface {
	Create(context.Context, dto.CreateDTO, dto.EnrichmentDTO) (int, error)

	Get(context.Context, dto.GetDTO, dto.FilterDTO, dto.SortDTO) ([]dto.User, error)
	GetById(context.Context, int) (dto.User, error)

	Update(context.Context, dto.UpdateDTO) (int, error)

	Delete(context.Context, int) (int, error)
}

type serviceEnrichment interface {
	Agify(string) (int, error)
	Genderize(string) (string, error)
	Nationalize(string) (string, error)
}

type UseCase struct {
	enrichment serviceEnrichment
	service    serviceUser

	logger log.Logger
}

func New(
	enrichment serviceEnrichment,
	service serviceUser,
	logger log.Logger,
) UseCase {

	return UseCase{
		enrichment: enrichment,
		service:    service,
		logger:     logger.WithField("unit", "enrichment"),
	}
}

func (u UseCase) Create(
	ctx context.Context,
	data dto.CreateDTO,
) (int, error) {

	age, err := u.enrichment.Agify(data.Name)
	if err != nil {
		return 0, err
	}

	gender, err := u.enrichment.Genderize(data.Name)
	if err != nil {
		return 0, err
	}

	country, err := u.enrichment.Nationalize(data.Name)
	if err != nil {
		return 0, err
	}

	enrichment := dto.EnrichmentDTO{
		Age:     age,
		Gender:  gender,
		Country: country,
	}

	return u.service.Create(ctx, data, enrichment)
}

func (u UseCase) Get(
	ctx context.Context,
	data dto.GetDTO,
	filter dto.FilterDTO,
	sort dto.SortDTO,
) ([]dto.User, error) {

	return u.service.Get(ctx, data, filter, sort)
}

func (u UseCase) GetById(
	ctx context.Context,
	id int,
) (dto.User, error) {

	return u.service.GetById(ctx, id)
}

func (u UseCase) Update(
	ctx context.Context,
	data dto.UpdateDTO,
) (int, error) {

	return u.service.Update(ctx, data)
}

func (u UseCase) Delete(
	ctx context.Context,
	id int,
) (int, error) {

	return u.service.Delete(ctx, id)
}
