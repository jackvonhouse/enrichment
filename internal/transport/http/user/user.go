package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackvonhouse/enrichment/internal/dto"
	"github.com/jackvonhouse/enrichment/internal/transport"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type useCaseUser interface {
	Create(context.Context, dto.CreateDTO) (int, error)

	Get(context.Context, dto.GetDTO, dto.FilterDTO, dto.SortDTO) ([]dto.User, error)
	GetById(context.Context, int) (dto.User, error)

	Update(context.Context, dto.UpdateDTO) (int, error)

	Delete(context.Context, int) (int, error)
}

type Transport struct {
	useCase useCaseUser

	logger log.Logger
}

func New(
	useCase useCaseUser,
	logger log.Logger,
) Transport {

	return Transport{
		useCase: useCase,
		logger:  logger.WithField("transport_type", "http"),
	}
}

func (t Transport) Handle(
	router *mux.Router,
) {

	router.HandleFunc("", t.Create).
		Methods(http.MethodPost)

	router.HandleFunc("", t.Get).
		Methods(http.MethodGet)

	router.HandleFunc("/{id:[0-9]+}", t.GetById).
		Methods(http.MethodGet)

	router.HandleFunc("/{id:[0-9]+}", t.Update).
		Methods(http.MethodPut)

	router.HandleFunc("/{id:[0-9]+}", t.Delete).
		Methods(http.MethodDelete)
}

func (t Transport) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	data := dto.CreateDTO{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		transport.Error(w, http.StatusInternalServerError, "invalid json structure")

		return
	}

	if data.Name == "" {
		transport.Error(w, http.StatusBadRequest, "empty name")

		return
	}

	if data.Surname == "" {
		transport.Error(w, http.StatusBadRequest, "empty surname")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := t.useCase.Create(ctx, data)
	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, map[string]any{"id": id})
}

func (t Transport) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	queries := r.URL.Query()

	name := queries.Get("name")
	surname := queries.Get("surname")
	patronymic := queries.Get("patronymic")

	age, err := transport.StringToInt(queries.Get("age"))
	if err != nil || age < 0 {
		age = 0
	}

	ageSort := queries.Get("age_sort_operator")

	genders := queries["gender"]
	countries := queries["country"]

	filter := dto.FilterDTO{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        age,
		AgeSort:    ageSort,
		Gender:     genders,
		Country:    countries,
	}

	sortBy := queries.Get("sort_by")
	if !transport.IsSortField(sortBy) {
		sortBy = "id"
	}

	sortOrder := queries.Get("sort_order")
	if !transport.IsSortOrder(sortOrder) {
		sortOrder = "desc"
	}

	sort := dto.SortDTO{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	limit, err := transport.StringToInt(queries.Get("limit"))
	if err != nil || limit <= 0 {
		// В зависимости от логики выбрасывать ошибку
		// или устанавливать limit по умолчанию
		// transport.Error(w, http.StatusBadRequest, "invalid limit")
		// return

		limit = 10
	}

	offset, err := transport.StringToInt(queries.Get("offset"))
	if err != nil || offset < 0 {
		// В зависимости от логики выбрасывать ошибку
		// или устанавливать offset по умолчанию
		// transport.Error(w, http.StatusBadRequest, "invalid offset")
		// return

		offset = 0
	}

	data := dto.GetDTO{
		Limit:  limit,
		Offset: offset,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := t.useCase.Get(ctx, data, filter, sort)
	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, users)
}

func (t Transport) GetById(
	w http.ResponseWriter,
	r *http.Request,
) {

	vars := mux.Vars(r)

	userID, err := transport.StringToInt(vars["id"])
	if err != nil || userID <= 0 {
		transport.Error(w, http.StatusBadRequest, "invalid enrichment id")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := t.useCase.GetById(ctx, userID)
	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, user)
}

func (t Transport) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	vars := mux.Vars(r)

	userID, err := transport.StringToInt(vars["id"])
	if err != nil || userID <= 0 {
		transport.Error(w, http.StatusBadRequest, "invalid user id")

		return
	}

	data := dto.UpdateDTO{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		transport.Error(w, http.StatusInternalServerError, "invalid json structure")

		return
	}

	data.ID = userID

	if data.Name == "" {
		transport.Error(w, http.StatusBadRequest, "name is empty")

		return
	}

	if data.Surname == "" {
		transport.Error(w, http.StatusBadRequest, "surname is empty")

		return
	}

	if data.Age <= 0 {
		transport.Error(w, http.StatusBadRequest, "invalid age")

		return
	}

	if data.Country == "" {
		transport.Error(w, http.StatusBadRequest, "country is empty")

		return
	}

	if data.Gender == "" {
		transport.Error(w, http.StatusBadRequest, "gender is empty")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := t.useCase.Update(ctx, data)
	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, map[string]any{"id": id})
}

func (t Transport) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	vars := mux.Vars(r)

	userID, err := transport.StringToInt(vars["id"])
	if err != nil || userID <= 0 {
		transport.Error(w, http.StatusBadRequest, "invalid user id")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := t.useCase.Delete(ctx, userID)
	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, map[string]any{"id": id})
}
