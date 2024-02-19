package user

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	pgerr "github.com/jackc/pgerrcode"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"github.com/jackvonhouse/enrichment/internal/errors"
	errpkg "github.com/jackvonhouse/enrichment/pkg/errors"
	"github.com/jackvonhouse/enrichment/pkg/log"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db     *sqlx.DB
	logger log.Logger
}

func New(
	db *sqlx.DB,
	logger log.Logger,
) Repository {

	return Repository{
		db:     db,
		logger: logger.WithField("unit", "enrichment"),
	}
}

func (r Repository) Create(
	ctx context.Context,
	create dto.CreateDTO,
	enrichment dto.EnrichmentDTO,
) (int, error) {

	query, args, err := sq.
		Insert("users").
		Columns(
			"name", "surname", "patronymic",
			"age", "gender", "country",
		).
		Values(
			create.Name, create.Surname, create.Patronymic,
			enrichment.Age, enrichment.Gender, enrichment.Country,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"name":       create.Name,
				"surname":    create.Surname,
				"patronymic": create.Patronymic,
				"age":        enrichment.Age,
				"gender":     enrichment.Gender,
				"country":    enrichment.Country,
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return 0, err
	}

	var userID int

	if err := r.db.GetContext(ctx, &userID, query, args...); err != nil {
		logger.Warnf("error on insert user: %s", err)

		if e, ok := err.(*pq.Error); ok {
			switch e.Code {

			case pgerr.UniqueViolation:
				return 0, errors.
					ErrAlreadyExists.
					New("user already exists").
					Wrap(err)

			case pgerr.ForeignKeyViolation:
				return 0, errors.
					ErrNotFound.
					New("user not found").
					Wrap(err)

			default:
				return 0, errors.
					ErrInternal.
					New("internal error").
					Wrap(err)
			}
		}

	}

	return userID, nil
}

func (r Repository) Get(
	ctx context.Context,
	get dto.GetDTO,
	filter dto.FilterDTO,
	sort dto.SortDTO,
) ([]dto.User, error) {

	sb := sq.
		Select(
			"id",
			"name", "surname", "patronymic",
			"age", "gender", "country",
		).
		From("users").
		OrderBy(
			fmt.Sprintf("%s %s", sort.SortBy, sort.SortOrder),
		).
		Offset(uint64(get.Offset)).
		Limit(uint64(get.Limit)).
		PlaceholderFormat(sq.Dollar)

	sb = r.where(sb, filter)

	query, args, err := sb.ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"offset": get.Offset,
				"limit":  get.Limit,
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return []dto.User{}, err
	}

	users := make([]dto.User, 0)

	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		logger.Warnf("error on get users: %s", err)

		if !errpkg.Is(err, sql.ErrNoRows) {
			return []dto.User{}, errors.
				ErrInternal.
				New("error on get users").
				Wrap(err)
		}

		return []dto.User{}, errors.
			ErrNotFound.
			New("haven't users").
			Wrap(err)
	}

	return users, nil
}

func (r Repository) GetById(
	ctx context.Context,
	id int,
) (dto.User, error) {

	query, args, err := sq.
		Select(
			"id",
			"name", "surname", "patronymic",
			"age", "gender", "country",
		).
		From("users").
		OrderBy("id").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"id": id,
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return dto.User{}, err
	}

	var user dto.User

	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		logger.Warnf("error on get user: %s", err)

		if !errpkg.Is(err, sql.ErrNoRows) {
			return dto.User{}, errors.
				ErrInternal.
				New("error on get users").
				Wrap(err)
		}

		return dto.User{}, errors.
			ErrNotFound.
			New("user not found").
			Wrap(err)
	}

	return user, nil
}

func (r Repository) Update(
	ctx context.Context,
	data dto.UpdateDTO,
) (int, error) {

	query, args, err := sq.
		Update("users").
		SetMap(map[string]interface{}{
			"name":       data.Name,
			"surname":    data.Surname,
			"patronymic": data.Patronymic,
			"age":        data.Age,
			"gender":     data.Gender,
			"country":    data.Country,
		}).
		Where(sq.Eq{"id": data.ID}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"id":         data.ID,
				"name":       data.Name,
				"surname":    data.Surname,
				"patronymic": data.Patronymic,
				"age":        data.Age,
				"gender":     data.Gender,
				"country":    data.Country,
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return 0, err
	}

	var userID int

	if err := r.db.GetContext(ctx, &userID, query, args...); err != nil {
		logger.Warnf("error on update user: %s", err)

		if e, ok := err.(*pq.Error); ok {
			switch e.Code {

			case pgerr.UniqueViolation:
				return 0, errors.
					ErrAlreadyExists.
					New("user already exists").
					Wrap(err)

			case pgerr.ForeignKeyViolation:
				return 0, errors.
					ErrNotFound.
					New("user not found").
					Wrap(err)

			default:
				return 0, errors.
					ErrInternal.
					New("internal error").
					Wrap(err)
			}
		}
	}

	return userID, err
}

func (r Repository) Delete(
	ctx context.Context,
	id int,
) (int, error) {

	query, args, err := sq.
		Delete("users").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"id": id,
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return 0, err
	}

	var userID int

	if err := r.db.GetContext(ctx, &userID, query, args...); err != nil {
		logger.Warnf("error on delete user: %s", err)

		if !errpkg.Is(err, sql.ErrNoRows) {
			return 0, errors.
				ErrInternal.
				New("error on delete user").
				Wrap(err)
		}

		return 0, errors.
			ErrNotFound.
			New("user not found").
			Wrap(err)
	}

	return userID, err
}
