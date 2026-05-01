package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"ticket-system/internal/config"
	"ticket-system/internal/model"
	pb "ticket-system/proto"

	"go.mongodb.org/mongo-driver/bson"
)

type InventoryService struct {
	pb.UnimplementedSeatServiceServer  // it will connect .proto file with this struct
}

func NewInventoryService() *InventoryService {
	collection := config.DB.Collection("seats")

	count , _ := collection.CountDocuments(context.Background(),bson.M{})

	if count ==0{

	for i := 1; i <= 10; i++ {
		seat := &model.Seat{
			SeatID: fmt.Sprintf("A%d",i),
			Status: "AVAILABLE",
		}
		collection.InsertOne(context.Background(), seat)
	}
}

	return &InventoryService{}
}

func (s *InventoryService) SeatStream(stream pb.SeatService_SeatStreamServer) error {
	for {
		event, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Println("Received:", event)

		if event.Action == "LOCK" {
			success := s.lockSeat(event.SeatId, event.UserId)

			if success {
				stream.Send(&pb.SeatEvent{
					SeatId: event.SeatId,
					Action: "LOCKED",
				})
			} else {
				stream.Send(&pb.SeatEvent{
					SeatId: event.SeatId,
					Action: "FAILED",
				})
			}
		
		}
	}
}

func (s *InventoryService) lockSeat(seatID, userID string) bool {

	key := "seatId:"+seatID
	 ok , err := config.ClientRadis.SetNX(config.Ctx, key , userID, 2*time.Minute).Result()
	 if err!=nil{
		log.Println("Redis error", err)
		return false 
	 } 
	 if ok {
		log.Println("Seat locked in Redis ", seatID)
		return true
	 }

	 log.Println("Seat already locked")
	 return false 


}

func (s *InventoryService) ReleaseSeat(ctx context.Context, req *pb.SeatRequest) (*pb.Empty, error) {
	key := "seatId:"+ req.SeatId
	_ , err := config.ClientRadis.Del(config.Ctx,key).Result()

	if err!=nil{
		log.Println("Redis Error")
		return &pb.Empty{}, err
	}
	log.Println("Seat released successfully", req.SeatId)

	return &pb.Empty{}, nil
}

func (s *InventoryService) ConfirmSeat(ctx context.Context, req *pb.SeatRequest) (*pb.Empty, error) {

	key := "seatId:" + req.SeatId

	// 🔥 STEP 1: CHECK LOCK EXISTS
	val, err := config.ClientRadis.Get(ctx, key).Result()
	if err != nil {
		log.Println("❌ Lock expired → payment invalid:", req.SeatId)
		return nil, fmt.Errorf("lock expired")
	}

	log.Println("Lock exists, owned by:", val)

	// 🔥 STEP 2: DELETE LOCK (consume it)
	_, err = config.ClientRadis.Del(ctx, key).Result()
	if err != nil {
		log.Println("Redis delete error:", err)
		return &pb.Empty{}, err
	}

	collection := config.DB.Collection("seats")

	filter := bson.M{
		"seatId": req.SeatId,
	}

	update := bson.M{
		"$set": bson.M{
			"status": "SOLD",
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error:", err)
		return nil, fmt.Errorf("lock expired")
	}

	if result.ModifiedCount == 1 {
		log.Println("✅ Seat sold:", req.SeatId)
			return &pb.Empty{}, nil


	} else {
		log.Println("❌ Seat not locked or already sold:", req.SeatId)
			return nil, fmt.Errorf("seat not updated")

	}

}
func (s *InventoryService) GetAvailableSeats(ctx context.Context, _ *pb.Empty) (*pb.SeatList, error) {

	collection := config.DB.Collection("seats")

	cursor, err := collection.Find(ctx, bson.M{"status": "AVAILABLE"})
	if err != nil {
		return nil, err
	}

	var result []*pb.Seat

	for cursor.Next(ctx) {
		var seat model.Seat
		cursor.Decode(&seat)
		
		key := "seatId:"+ seat.SeatID
		exists , err := config.ClientRadis.Exists(config.Ctx,key).Result()
		if err!=nil{
			log.Println("Redis error",err)
		}

		if exists==0{
		result = append(result, &pb.Seat{
			SeatId: seat.SeatID,
			Status: seat.Status,
			UserId: seat.UserID,
		})
	}
	}

	return &pb.SeatList{Seats: result}, nil
}