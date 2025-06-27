package gqlhandler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/axmz/go-port-service/internal/services/port"
	graphql "github.com/axmz/go-port-service/internal/transport/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQLHandler = handler.Server

func InitGql(portSvc *port.Service) *GraphQLHandler {
	gqlsrv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{
		PortService: portSvc,
	}}))
	gqlsrv.AddTransport(transport.Options{})
	gqlsrv.AddTransport(transport.GET{})
	gqlsrv.AddTransport(transport.POST{})
	gqlsrv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	gqlsrv.Use(extension.Introspection{})
	gqlsrv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})
	return gqlsrv
}
