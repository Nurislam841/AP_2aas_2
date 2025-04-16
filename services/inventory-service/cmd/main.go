package main

import (
	"database/sql"
	"fmt"
	"log"

	// "github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	grpc_hand "inventory-service/delivery/grpc"
	// "inventory-service/delivery/http"
	"inventory-service/internal/repository"
	"inventory-service/internal/usecase"
	"net"
	"proto/inventorypb"
)

func main() {

	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost:5432/inventory_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Println("Connected to the database")

	repo := repository.NewProductRepository(db)
	usecase := usecase.NewProductUsecase(repo)

	// router := gin.Default()
	// http.NewProductHandler(router, usecase)

	// router.Run(":5001")

	grpcServer := grpc.NewServer()

	userServiceServer := grpc_hand.NewInventoryServer(usecase)

	inventorypb.RegisterInventoryServiceServer(grpcServer, userServiceServer)

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Inventary Service running on :5001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
