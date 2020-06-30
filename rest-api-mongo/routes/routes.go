package routes

import (
	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
	"net/http"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/configs"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/db"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/handlers"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/handlers/social_logins"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/middlewares"
	"github.com/gorilla/mux"
)

// SetupRoutes , sets up all the routes and returns router instance
func SetupRoutes(repos *db.Repositories, serverConfigs *configs.Config, mss protos.NotificationClient) *mux.Router {
	router := mux.NewRouter()

	// Handlers
	todoHandler := handlers.NewTodoHandler(repos.TodoRepo, mss)
	userHandler := handlers.NewUserHandler(repos.UserRepo, serverConfigs, mss)
	googleHandler := social_logins.NewGoogleHandler(serverConfigs.Google)
	jwtMiddleWare := middlewares.NewJWTMiddleWare(serverConfigs.JWT)

	router.HandleFunc("/login", userHandler.PerformLogin).Methods(http.MethodPost)
	router.HandleFunc("/signup", userHandler.NewUserSignUp).Methods(http.MethodPost)
	router.HandleFunc("/resetpassword", userHandler.RequestResetPassword).Methods(http.MethodPost)
	router.HandleFunc("/resetpassword/{slug}", userHandler.ResetPassword).Methods(http.MethodGet)
	router.HandleFunc("/verify/{slug}", userHandler.VerifyEmail).Methods(http.MethodGet)

	// todoRouter, handles all the request coming on /todo path
	todoRouter := router.PathPrefix("/todo").Subrouter()
	todoRouter.Use(jwtMiddleWare.ValidateAndRefreshToken)
	todoRouter.HandleFunc("/items", todoHandler.GetAllTodos).Methods(http.MethodGet)
	todoRouter.HandleFunc("/items/{id:[0-9]+}", todoHandler.GetTodo).Methods(http.MethodGet)
	todoRouter.HandleFunc("/items", todoHandler.AddTodo).Methods(http.MethodPost)
	todoRouter.HandleFunc("/items/{id:[0-9]+}", todoHandler.DeleteTodo).Methods(http.MethodDelete)

	// userRouter, handlers all the request coming on /user path
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(jwtMiddleWare.ValidateAndRefreshToken)
	userRouter.HandleFunc("", userHandler.SaveUser).Methods(http.MethodPost)
	userRouter.HandleFunc("", userHandler.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/verify", userHandler.RequestVerifyEmail).Methods(http.MethodGet)

	// googleRouter, handles all request coming on /google
	googleRouter := router.PathPrefix("/google").Subrouter()
	googleRouter.HandleFunc("", googleHandler.GoogleHome)
	googleRouter.HandleFunc("/login", googleHandler.GoogleLogin)
	googleRouter.HandleFunc("/callback", googleHandler.GoogleCallback)

	return router
}
