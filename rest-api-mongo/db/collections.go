package db

import "github.com/frost060/go-microservice-basic/rest-api-mongo/logging"

const (
	TODO = "todos"
	USER = "users"
)

// Repositories , holds instances of all the collection
type Repositories struct {
	UserRepo *UserRepo
	TodoRepo *TodoRepo
}

// SetupRepositories , connects to DB and return instance of all the collections
func SetupRepositories(log *logging.LogWrapper) *Repositories {
	db := ConnectDB(log)

	// Repositories
	todoRepo := NewTodoRepo(db, log)
	userRepo := NewUserRepo(db, log)

	return &Repositories{
		UserRepo: userRepo,
		TodoRepo: todoRepo,
	}
}
