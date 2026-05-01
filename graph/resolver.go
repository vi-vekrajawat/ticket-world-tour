package graph

import (
	"log"
	"ticket-system/graph/model"

	pb "ticket-system/proto"
	"google.golang.org/grpc"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct{
	SeatClient pb.SeatServiceClient
	SeatChannel chan *model.Seat
}


func NewResolver() *Resolver {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to connect gRPC:", err)
	}

	client := pb.NewSeatServiceClient(conn)

	return &Resolver{
		SeatClient:  client,
		SeatChannel: make(chan *model.Seat, 10),
	}
}
