package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB(){
	client , err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err!=nil{
		panic("database connenction failed ")
	}

	ctx, cancel := context.WithTimeout(context.Background(),time.Second*10)
	defer cancel()

err = client.Connect(ctx)
	if err!=nil{
		log.Fatal("database connection failed ", err)	
	
}

	DB = client.Database("ticket_system")

	log.Println("Database connected successfully")
}