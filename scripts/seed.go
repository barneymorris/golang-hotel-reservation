package main

import (
	"context"
	"log"

	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/betelgeusexru/golang-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	roomStore db.RoomStore
	hotelStore db.HotelStore
	ctx = context.Background()
)

func seedHotel(name, location string) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Type: types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type: types.DeluxeRoomType,
			BasePrice: 199.9,
		},
		{
			Type: types.SeaSideRoomType,
			BasePrice: 122.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	
	}
}

func main() {
	seedHotel("Bellucia", "France")
	seedHotel("The cozy hotel", "The Nederlands")
	seedHotel("Dont die in your sleep", "London")
}

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}