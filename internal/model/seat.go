package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SeatStatus string

const (
	Available SeatStatus = "AVAILABLE"
	Locked    SeatStatus = "LOCKED"
	Sold      SeatStatus = "SOLD"
)

type Seat struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	SeatID string 		   `bson:"seatId"`
	Status string  `bson:"status"`
	UserID string       `bson:"userId"`
}