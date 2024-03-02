package transport

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gorilla/mux"
	"github.com/jackvonhouse/enrichment/app/usecase"
	"github.com/jackvonhouse/enrichment/internal/transport/graphql"
	graphqlUser "github.com/jackvonhouse/enrichment/internal/transport/graphql/user"
	httpUser "github.com/jackvonhouse/enrichment/internal/transport/http/user"
	"github.com/jackvonhouse/enrichment/internal/transport/router"
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
		"/user": httpUser.New(useCase.User, transportLogger),
	})

	h := graphqlUser.New(useCase.User, transportLogger)

	srv := handler.NewDefaultServer(
		graphql.NewExecutableSchema(
			graphql.Config{
				Resolvers: &h,
			},
		),
	)

	r.Router().Handle("/graphql/user", srv)

	return Transport{
		router: r,
	}
}

func (t Transport) Router() *mux.Router { return t.router.Router() }
