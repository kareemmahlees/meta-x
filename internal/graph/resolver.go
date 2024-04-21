package graph

import (
	"github.com/kareemmahlees/meta-x/internal/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Storage db.Storage
}
