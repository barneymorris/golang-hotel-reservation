package testutils

import (
	"context"
	"log"
	"testing"

	"github.com/betelgeusexru/golang-hotel-reservation/config"
	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Testmongouri = "mongodb://localhost:27017" 

type Testdb struct {
	db.UserStore
}

func (tdb *Testdb) Teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func SetupDatabase(t *testing.T) *Testdb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(Testmongouri))
	if err != nil {
		log.Fatal(err)
	}

	return &Testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func SetupFiberApp() *fiber.App {
	app := fiber.New(config.Config)
	return app
}