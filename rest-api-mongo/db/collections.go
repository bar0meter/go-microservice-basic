package db

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
func SetupRepositories() *Repositories {
	db := ConnectDB()

	// Repositories
	todoRepo := NewTodoRepo(db)
	userRepo := NewUserRepo(db)

	return &Repositories{
		UserRepo: userRepo,
		TodoRepo: todoRepo,
	}
}
