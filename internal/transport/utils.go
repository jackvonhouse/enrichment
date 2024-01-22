package transport

import (
	"net/http"
	"strconv"

	"github.com/jackvonhouse/enrichment/internal/errors"
	errpkg "github.com/jackvonhouse/enrichment/pkg/errors"
)

func StringToInt(valueStr string) (int, error) {
	valueInt, err := strconv.Atoi(valueStr)

	if err != nil {
		return 0, err
	}

	return valueInt, nil
}

var DefaultErrorHttpCodes = map[uint32]int{
	errors.ErrInternal.TypeId:       http.StatusInternalServerError,
	errors.ErrCantEnrichment.TypeId: http.StatusInternalServerError,
	errors.ErrAlreadyExists.TypeId:  http.StatusConflict,
	errors.ErrNotFound.TypeId:       http.StatusNotFound,
}

func ErrorToHttpResponse(
	err error, codes map[uint32]int,
) (int, string) {

	if errpkg.Has(err, errors.ErrInternal) {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)
	}

	code := codes[errpkg.TypeId(err)]

	if code == 0 {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)
	}

	return code, err.Error()
}

var defaultSortFields = map[string]bool{
	"id":         true,
	"name":       true,
	"surname":    true,
	"patronymic": true,
	"age":        true,
	"gender":     true,
	"country":    true,
}

var defaultSortOrders = map[string]bool{
	"asc":  true,
	"desc": true,
}

func IsSortField(value string) bool {
	_, ok := defaultSortFields[value]
	return ok
}

func IsSortOrder(value string) bool {
	_, ok := defaultSortOrders[value]
	return ok
}
