package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/betelgeusexru/golang-hotel-reservation/api"
	"github.com/betelgeusexru/golang-hotel-reservation/api/middleware"
	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError

		var e *fiber.Error
        if errors.As(err, &e) {
            code = e.Code
        }

		ctx.Status(code)
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	store := &db.Store{
		Hotel: hotelStore,
		Room: roomStore,
		User: userStore,
	}

	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(userStore)

	app := fiber.New(config)

	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication)

	auth.Post("/auth", authHandler.HandleAuthenticate)

	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)  

	app.Listen(*listenAddr);
}