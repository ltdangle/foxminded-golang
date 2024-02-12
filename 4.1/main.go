package main

import (
	"log"
	"net"

	grpcSrv "grpc4_1/pkg/grpc"
	"grpc4_1/pkg/model"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config.
	type config struct {
		NetPort  string `env:"NET_PORT" validate:"required"`
		MysqlDsn string `env:"MYSQL_DSN" validate:"required"`
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	// Validate config.
	validator := validator.New()
	err = validator.Struct(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := sqlx.Open("mysql", cfg.MysqlDsn)
	if err != nil {
		log.Fatal("failed to connect database")
	}

	userRepo := model.NewUserRepo(db)

	// Specify the port on which the server should listen
	netAddr := "localhost:" + cfg.NetPort

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", netAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	grpcSrv.RegisterUserServiceServer(grpcServer, grpcSrv.NewUserService(userRepo))

	// Start the server
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
