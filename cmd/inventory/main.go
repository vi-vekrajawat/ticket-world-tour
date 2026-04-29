package main

import (
	"log"
	"net"
	"ticket-system/internal/config"
	service "ticket-system/internal/service/inventory-service"

	pb "ticket-system/proto"

	"google.golang.org/grpc"
)

func main() {

	config.ConnectDB()
	config.ConnectRadis()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	inventoryService := service.NewInventoryService()
	pb.RegisterSeatServiceServer(grpcServer, inventoryService)

	log.Println("Inventory Service running on 50051")

	grpcServer.Serve(lis)
}