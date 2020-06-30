package main

import (
	"fmt"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/db"
	log "github.com/frost060/go-microservice-basic/basic-messaging-service/logging"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/server"
	"google.golang.org/grpc"
)

func main() {
	log.Info("Starting notification service...")

	log.Info("Loading configs from .env file")
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
		os.Exit(1)
	}

	serverConfig := configs.NewConfig()

	gs := grpc.NewServer()
	log.Info("Created new grpc server...")

	redis := db.NewRedisClient(serverConfig)

	ms := server.NewMessageService(serverConfig, redis)
	log.Info("Create new message service...")

	protos.RegisterNotificationServer(gs, ms)
	log.Info("Successfully registered notification service")

	reflection.Register(gs)

	log.Info("Notification service running on port: 9092")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	_ = gs.Serve(l)
}
