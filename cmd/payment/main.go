package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"ticket-system/internal/worker"
	pb "ticket-system/proto"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewSeatServiceClient(conn)

	stream, err := client.SeatStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	

	go worker.StartWorker(client)

	res , err := client.GetAvailableSeats(context.Background(),&pb.Empty{})
	if err!=nil {
		log.Println("Something went wrong while fetching all the available seats")
	}

	for _ , seat := range res.Seats{
		log.Printf("Seat : %s Status %s\n", seat.SeatId , seat.Status)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the seat Id like (A1 , A2 , A3 .......)")
	seatId , err := reader.ReadString('\n')
	seatId = seatId[:len(seatId)-1] 
	seatId = strings.TrimSpace(seatId)
	if err != nil {
		log.Fatal(err)
	}
	


	// simulate user selecting seat
	event := &pb.SeatEvent{
		SeatId: seatId,
		UserId: "user2",
		Action: "LOCK",
	}

	stream.Send(event)

	go func() {
	for {
		res, err := stream.Recv()
		if err != nil {
			log.Println("Stream error:", err)
			return
		}

		log.Println("Response from inventory:", res)

		if res.Action == "LOCKED" {
			log.Println("✅ Seat locked, starting payment:", res.SeatId)
			worker.PaymentQueue <- res
		} else {
			log.Println("❌ Seat not available:", res.SeatId)
		}
	}
}()

	select {}
}