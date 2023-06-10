package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/betelgeusexru/golang-hotel-reservation/api"
	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/betelgeusexru/golang-hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	ctx context.Context
)

func main() {
	hotelStore := db.NewMongoHotelStore(client)

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store := db.Store{
		User: db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room: db.NewMongoRoomStore(client, hotelStore),
		Hotel: db.NewMongoHotelStore(client),
	}

	user := fixtures.AddUser(&store, "james", "foo", false)
	fmt.Println("james ->", api.CreateTokenFromUser(user))

	fixtures.AddUser(&store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(user))

	hotel := fixtures.AddHotel(&store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(&store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(&store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))

	fmt.Println("booking ->", booking.ID)
}

func init() {
	ctx = context.Background()

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
}