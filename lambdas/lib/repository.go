// Package lib provides generic repository interfaces and types for database operations.
package lib

import (
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/lib/conversation"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/status"
	"github.com/kmesiab/equilibria/lambdas/lib/user"
)

// RepositoryType is the interface for all repositories. Implement this interface
// when creating a new repository type.
type RepositoryType interface {
	conversation.ConversationRepository | message.Repository | status.Repository | user.UserRepository
}

// Repository is a generic interface for database operations on type T.
// It includes methods for common CRUD operations and allows setting a specific
// database instance, which can be a real or mock database.
type Repository[T any] interface {
	// SetDBInstance sets the database instance for the repository.
	SetDBInstance(db *gorm.DB)

	// FindByID retrieves an entity of type T by its ID from the database.
	FindByID(id int64) (*T, error)

	// Create inserts a new entity of type T into the database.
	Create(t *T) error

	// Update modifies an existing entity of type T in the database.
	Update(t *T) error

	// Delete performs a soft delete on an entity of type T in the database by its ID.
	Delete(id int64) error
}
