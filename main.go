package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/betelgeusexru/golang-hotel-reservation/api"
	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

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

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	mongoStore := db.NewMongoUserStore(client, db.DBNAME)
	userHandler := api.NewUserHandler(mongoStore)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	
	app.Listen(*listenAddr);
}