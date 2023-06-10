package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/betelgeusexru/golang-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, uid,rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID: uid,
		RoomID: rid,
		FromDate: from,
		TillDate: till,
	}

	_, err := store.Booking.InsertBooking(context.Background(), booking)

	if err != nil {
		log.Fatal(err)
	}

	return booking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size: size,
		Seaside: ss,
		Price: price,
		HotelID: hid,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	
	hotel := types.Hotel{
		Name: name,
		Location: loc,
		Rooms: roomIDS,
		Rating: rating,
	}


	insertedHotel, err := store.Hotel.Insert(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddUser(store *db.Store,  fn, ln string, admin bool) *types.User { 
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName: ln,
		Password: fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}