package transport

import (
	"github.com/gorilla/mux"
	"github.com/jackvonhouse/enrichment/app/usecase"
	"github.com/jackvonhouse/enrichment/internal/transport/router"
	"github.com/jackvonhouse/enrichment/internal/transport/user"
	"github.com/jackvonhouse/enrichment/pkg/log"
)

type Transport struct {
	router *router.Router
}

func New(
	useCase usecase.UseCase,
	logger log.Logger,
) Transport {

	transportLogger := logger.WithField("layer", "transport")

	r := router.New("/api/v1")

	r.Handle(map[string]router.Handlify{
		"/user": user.New(useCase.User, transportLogger),
	})

	return Transport{
		router: r,
	}
}

func (t Transport) Router() *mux.Router { return t.router.Router() }
