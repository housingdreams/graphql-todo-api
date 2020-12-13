package graph

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/uuid"
	"github.com/leminhson2398/todo-api/internal/db"
	"github.com/leminhson2398/todo-api/internal/utils"
)

// NewHandler endpoint handles all graphql request sent to server
func NewHandler(repo db.Repository) http.Handler {
	config := Config{
		Resolvers: &Resolver{
			Repository: repo,
		},
	}

	srv := handler.New(NewExecutableSchema(config))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	if isProd := os.Getenv("PRODUCTION") == "true"; isProd {
		srv.Use(extension.FixedComplexityLimit(10))
	} else {
		srv.Use(extension.Introspection{})
	}
	return srv
}

// NewPlaygroundHandler display graphql playground
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}

// GetCurrentUserID returns current user's ID
func GetCurrentUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(utils.UserIDKey).(uuid.UUID)
	return userID, ok
}
