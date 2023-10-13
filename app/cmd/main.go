package main

import (
	"flag"
	"github.com/zumosik/grpc-user-auth-service-go/config"
	"github.com/zumosik/grpc-user-auth-service-go/pb"
	"github.com/zumosik/grpc-user-auth-service-go/server"
	"github.com/zumosik/grpc-user-auth-service-go/storage/postgres"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	cfgPath string
)

func init() {
	flag.StringVar(&cfgPath, "cfg-path", "configs/cfg.yml", "Path to your yaml config")
}

func main() {
	log.Println("---- User service ----") // very pretty logs
	cfg := config.MustConfig(cfgPath)
	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}
	defer func(storage *postgres.Storage) {
		_ = storage.Close()
	}(storage)
	log.Println("Connected to db!")

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("can't create lis: %v", err)
	}

	srv := server.New(storage, time.Duration(cfg.Timeout)*time.Second)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, srv)

	log.Println("Starting server!")

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
