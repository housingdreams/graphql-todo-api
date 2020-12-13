//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"sync"

	"github.com/leminhson2398/todo-api/internal/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository db.Repository
	mu         sync.Mutex
}
