package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	grpc_hand "order-service/delivery/grpc"
	// "github.com/gin-gonic/gin"
	// "order-service/delivery/http"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"proto/orderpb"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Println("Connected to the database")

	repo := repository.NewOrderRepository(db)
	usecase := usecase.NewOrderUsecase(repo)

	// router := gin.Default()
	// http.NewOrderHandler(router, usecase)

	// router.Run(":5002")
	grpcServer := grpc.NewServer()

	userServiceServer := grpc_hand.NewOrderServer(usecase)

	orderpb.RegisterOrderServiceServer(grpcServer, userServiceServer)

	lis, err := net.Listen("tcp", ":5002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Order Service running on :5002")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
