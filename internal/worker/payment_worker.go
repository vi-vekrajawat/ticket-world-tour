package worker

import (
	"context"
	"log"
	"time"

	pb "ticket-system/proto"
)

var PaymentQueue = make(chan *pb.SeatEvent, 100)

func StartWorker(client pb.SeatServiceClient) {
	for event := range PaymentQueue {
		go handleSeat(event, client)
	}
}

func handleSeat(event *pb.SeatEvent, client pb.SeatServiceClient) {
	log.Println("Timer started for:", event.SeatId)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	

	<-fakePaymentSuccess()
		// log.Println("Payment success:", event.SeatId)
		res , err := client.ConfirmSeat(ctx, &pb.SeatRequest{
			SeatId: event.SeatId,
		})
		if err!=nil{
			log.Println("Confirm seat grpc err")
			return
		} else{
			log.Println("Confirmed the seat",res)
		}

	
}

func fakePaymentSuccess() <-chan bool {
	ch := make(chan bool)

	go func() {
		time.Sleep(18 * time.Second)
		ch <- true
	}()

	return ch
}