package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/configs"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/logging"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/routes"
	"github.com/joho/godotenv"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/db"
	gohandlers "github.com/gorilla/handlers"

	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
)

func main() {
	log := logging.NewLogger()
	log.Info("Starting the application...")

	repos := db.SetupRepositories(log)

	log.Info("Loading configs from .env file")
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
		os.Exit(1)
	}

	serverConfigs := configs.NewConfig()
	log.Info("Successfully loaded configs")

	// Message Service Client
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		log.Error("Error while connecting to messaging service")
		os.Exit(1)
	}

	defer conn.Close()

	mssClient := protos.NewNotificationClient(conn)

	router := routes.SetupRoutes(repos, serverConfigs, mssClient, log)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      ch(router),        // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Info("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			log.Error("Unable to start server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Info("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(ctx)
}
