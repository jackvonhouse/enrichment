package user

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"strings"
)

func (r Repository) where(
	builder sq.SelectBuilder,
	filter dto.FilterDTO,
) sq.SelectBuilder {

	builder = r.whereName(builder, filter.Name)
	builder = r.whereSurname(builder, filter.Surname)
	builder = r.wherePatronymic(builder, filter.Patronymic)
	builder = r.whereAge(builder, filter.Age, filter.AgeSort)
	builder = r.whereGender(builder, filter.Gender)
	builder = r.whereCountry(builder, filter.Country)

	return builder
}

func (r Repository) whereName(
	builder sq.SelectBuilder,
	name string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(name)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"name": fmt.Sprintf("%%%s%%", strings.ToLower(name)),
		},
	)
}
func (r Repository) whereSurname(
	builder sq.SelectBuilder,
	surname string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(surname)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"surname": fmt.Sprintf("%%%s%%", strings.ToLower(surname)),
		},
	)
}

func (r Repository) wherePatronymic(
	builder sq.SelectBuilder,
	patronymic string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(patronymic)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"patronymic": fmt.Sprintf("%%%s%%", strings.ToLower(patronymic)),
		},
	)
}

func (r Repository) whereAge(
	builder sq.SelectBuilder,
	age int,
	sortParam string,
) sq.SelectBuilder {

	if age <= 0 {
		return builder
	}

	switch sortParam {
	case "eq":
		return builder.Where(sq.Eq{"age": age})
	case "ne":
		return builder.Where(sq.NotEq{"age": age})
	case "gt":
		return builder.Where(sq.Gt{"age": age})
	case "ge":
		return builder.Where(sq.GtOrEq{"age": age})
	case "lt":
		return builder.Where(sq.Lt{"age": age})
	case "le":
		return builder.Where(sq.LtOrEq{"age": age})
	default:
		return builder
	}
}

func (r Repository) whereGender(
	builder sq.SelectBuilder,
	genders []string,
) sq.SelectBuilder {

	if len(genders) == 0 {
		return builder
	}

	orCondition := sq.Or{}
	for _, gender := range genders {
		orCondition = append(
			orCondition,
			sq.Eq{"gender": strings.ToLower(gender)},
		)
	}

	return builder.Where(orCondition)
}

func (r Repository) whereCountry(
	builder sq.SelectBuilder,
	countries []string,
) sq.SelectBuilder {

	if len(countries) == 0 {
		return builder
	}

	orCondition := sq.Or{}
	for _, country := range countries {
		orCondition = append(
			orCondition,
			sq.Eq{"country": strings.ToUpper(country)},
		)
	}

	return builder.Where(orCondition)
}
