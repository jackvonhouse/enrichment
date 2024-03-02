package user

import (
	"context"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"github.com/jackvonhouse/enrichment/internal/errors"
	graphql1 "github.com/jackvonhouse/enrichment/internal/transport/graphql"
	"github.com/jackvonhouse/enrichment/internal/transport/graphql/models"
)

var (
	ErrEmptyName     = errors.ErrEmptyField.New("empty name")
	ErrEmptySurname  = errors.ErrEmptyField.New("empty surname")
	ErrEmptyCountry  = errors.ErrEmptyField.New("empty country")
	ErrEmptyGender   = errors.ErrEmptyField.New("empty gender")
	ErrInvalidUserId = errors.ErrInvalidValue.New("invalid user id")
	ErrInvalidAge    = errors.ErrInvalidValue.New("invalid age")
)

func (r *mutationResolver) Create(
	ctx context.Context,
	input models.CreateInput,
) (int, error) {

	if input.Name == "" {
		r.logger.Warn("user name is empty")

		return 0, ErrEmptyName
	}

	if input.Surname == "" {
		r.logger.Warn("user surname is empty")

		return 0, ErrEmptySurname
	}

	var patronymic string

	if input.Patronymic != nil {
		patronymic = *input.Patronymic
	}

	dataInput := dto.CreateDTO{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: patronymic,
	}

	return r.useCase.Create(ctx, dataInput)
}

func (r *queryResolver) Get(
	ctx context.Context,
	get *models.GetInput,
	filter *models.FilterInput,
	sort *models.SortInput,
) ([]models.User, error) {

	var (
		getInput    dto.GetDTO
		filterInput dto.FilterDTO
		sortInput   dto.SortDTO
	)

	if get == nil {
		getInput = dto.GetDTO{
			Limit:  10,
			Offset: 0,
		}
	} else {
		var (
			limit  = 10
			offset = 0
		)

		if get.Limit != nil {
			limit = *get.Limit
		}

		if get.Offset != nil {
			offset = *get.Offset
		}

		getInput = dto.GetDTO{
			Limit:  limit,
			Offset: offset,
		}
	}

	if filter == nil {
		filterInput = dto.FilterDTO{
			Name:       "",
			Surname:    "",
			Patronymic: "",
			Age:        0,
			AgeSort:    "",
			Gender:     []string{},
			Country:    []string{},
		}
	} else {
		var (
			name       = ""
			surname    = ""
			patronymic = ""
			age        = 0
			ageSort    = ""
			gender     []string
			country    []string
		)

		if filter.Name != nil {
			name = *filter.Name
		}

		if filter.Surname != nil {
			surname = *filter.Surname
		}

		if filter.Patronymic != nil {
			patronymic = *filter.Patronymic
		}

		if filter.Age != nil {
			age = *filter.Age
		}

		if filter.AgeSort != nil {
			ageSort = *filter.AgeSort
		}

		if filter.Gender != nil {
			gender = make([]string, len(filter.Gender))
			copy(gender, filter.Gender)
		}

		if filter.Country != nil {
			country = make([]string, len(filter.Country))
			copy(country, filter.Country)
		}

		filterInput = dto.FilterDTO{
			Name:       name,
			Surname:    surname,
			Patronymic: patronymic,
			Age:        age,
			AgeSort:    ageSort,
			Gender:     gender,
			Country:    country,
		}
	}

	if sort == nil {
		sortInput = dto.SortDTO{
			SortBy:    "id",
			SortOrder: "desc",
		}
	} else {
		var (
			sortBy    = "id"
			sortOrder = "desc"
		)

		if sort.SortBy != nil {
			sortBy = *sort.SortBy
		}

		if sort.SortOrder != nil {
			sortOrder = *sort.SortOrder
		}

		sortInput = dto.SortDTO{
			SortBy:    sortBy,
			SortOrder: sortOrder,
		}
	}

	users, err := r.useCase.Get(ctx, getInput, filterInput, sortInput)
	if err != nil {
		r.logger.Warnf("error getting users: %s", err)

		return []models.User{}, err
	}

	usersModel := make([]models.User, len(users))

	for i, user := range users {
		user := user

		usersModel[i] = models.User{
			ID:         user.ID,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: &user.Patronymic,
			Age:        &user.Age,
			Country:    &user.Country,
			Gender:     &user.Gender,
		}
	}

	return usersModel, nil
}

func (r *queryResolver) GetByID(
	ctx context.Context,
	id int,
) (models.User, error) {

	if id <= 0 {
		r.logger.Warn("invalid user id")

		return models.User{}, ErrInvalidUserId
	}

	user, err := r.useCase.GetById(ctx, id)
	if err != nil {
		r.logger.Warn("error getting user by id: %s", err)

		return models.User{}, err
	}

	return models.User{
		ID:         user.ID,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: &user.Patronymic,
		Age:        &user.Age,
		Gender:     &user.Gender,
		Country:    &user.Country,
	}, nil
}

func (r *mutationResolver) Update(
	ctx context.Context,
	input models.UpdateInput,
) (int, error) {

	if input.ID <= 0 {
		r.logger.Warn("invalid user id")

		return 0, ErrInvalidUserId
	}

	if input.Name == "" {
		r.logger.Warn("empty user name")

		return 0, ErrEmptyName
	}

	if input.Surname == "" {
		r.logger.Warn("empty user surname")

		return 0, ErrEmptySurname
	}

	if input.Age <= 0 {
		r.logger.Warn("empty user age")

		return 0, ErrInvalidAge
	}

	if input.Country == "" {
		r.logger.Warn("empty user country")

		return 0, ErrEmptyCountry
	}

	if input.Gender == "" {
		r.logger.Warn("empty user gender")

		return 0, ErrEmptyGender
	}

	updateInput := dto.UpdateDTO{
		ID:         input.ID,
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: *input.Patronymic,
		Age:        input.Age,
		Country:    input.Country,
		Gender:     input.Gender,
	}

	return r.useCase.Update(ctx, updateInput)
}

func (r *mutationResolver) Delete(
	ctx context.Context,
	id int,
) (int, error) {

	if id <= 0 {
		r.logger.Warn("invalid user id")

		return 0, ErrInvalidUserId
	}

	return r.useCase.Delete(ctx, id)
}

func (r *Resolver) Mutation() graphql1.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
