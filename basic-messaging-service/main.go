package main

import (
	"fmt"
	"net"
	"os"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	configs := configs.NewConfig()

	gs := grpc.NewServer()

	ms := server.NewMessageService(log, configs)

	notifications.RegisterNotificationServer(gs, ms)

	reflection.Register(gs)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	gs.Serve(l)
}
