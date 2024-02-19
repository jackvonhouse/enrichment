package enrichment

import (
	"encoding/json"
	"fmt"
	"github.com/jackvonhouse/enrichment/internal/errors"
	"github.com/jackvonhouse/enrichment/pkg/log"
	"net/http"
)

var (
	ErrCantAgify       = errors.ErrCantEnrichment.New("can't agify")
	ErrCantGenderize   = errors.ErrCantEnrichment.New("can't genderize")
	ErrCantNationalize = errors.ErrCantEnrichment.New("can't nationalize")
)

type Service struct {
	logger log.Logger
}

func New(
	logger log.Logger,
) Service {

	return Service{
		logger: logger.WithField("unit", "enrichment"),
	}
}

func (s Service) Agify(name string) (int, error) {
	logger := s.logger.WithField("name", name)

	resp, err := http.Get(s.agifyUrl(name))
	if err != nil {
		logger.Warnf("can't get from agify: %s", err)

		return 0, ErrCantAgify
	}

	defer resp.Body.Close()

	var data struct {
		Count int    `json:"count"`
		Age   int    `json:"age"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Warnf("can't decode agify: %s", err)

		return 0, ErrCantAgify
	}

	return data.Age, nil
}

func (s Service) Genderize(name string) (string, error) {
	logger := s.logger.WithField("name", name)

	resp, err := http.Get(s.genderizeUrl(name))
	if err != nil {
		logger.Warnf("can't get from genderize: %s", err)

		return "", ErrCantGenderize
	}

	defer resp.Body.Close()

	var data struct {
		Count       int     `json:"count"`
		Name        string  `json:"name"`
		Gender      string  `json:"gender"`
		Probability float64 `json:"probability"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Warnf("can't decode genderize: %s", err)

		return "", ErrCantGenderize
	}

	return data.Gender, nil
}

func (s Service) Nationalize(name string) (string, error) {
	logger := s.logger.WithField("name", name)

	resp, err := http.Get(s.nationalizeUrl(name))
	if err != nil {
		logger.Warnf("can't get from nationalize: %s", err)

		return "", ErrCantNationalize
	}

	defer resp.Body.Close()

	var data struct {
		Count   int    `json:"count"`
		Name    string `json:"name"`
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Warnf("can't decode nationalize: %s", err)

		return "", ErrCantNationalize
	}

	if len(data.Country) == 0 {
		logger.Warnf("error on nationalize: countries length is %d", len(data.Country))

		return "", ErrCantNationalize
	}

	return data.Country[0].CountryID, nil
}

func (s Service) agifyUrl(name string) string {
	return fmt.Sprintf("https://api.agify.io?name=%s", name)
}

func (s Service) genderizeUrl(name string) string {
	return fmt.Sprintf("https://api.genderize.io?name=%s", name)
}

func (s Service) nationalizeUrl(name string) string {
	return fmt.Sprintf("https://api.nationalize.io?name=%s", name)
}
