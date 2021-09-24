package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Asymmetriq/url_shortener/internal/service"
	"github.com/Asymmetriq/url_shortener/pkg/api"
	gw "github.com/Asymmetriq/url_shortener/pkg/pb/api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// server
	srv := &service.Service{
		DB: db,
	}

	// GRPC
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	grpcOpts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(grpcOpts...)

	api.RegisterServiceServer(grpcServer, srv)

	log.Println("gRPC server listening on :5000")
	go grpcServer.Serve(listener)

	// HTTP
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	gwmux := runtime.NewServeMux()

	err = gw.RegisterServiceHandlerServer(ctx, gwmux, srv)
	if err != nil {
		log.Fatal(err)
	}

	// Create a gRPC Gateway server
	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: gwmux,
	}

	log.Println("gRPC server listening on :8000")
	httpServer.ListenAndServe()
}
