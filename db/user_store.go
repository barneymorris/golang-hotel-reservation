package db

import (
	"context"

	"github.com/betelgeusexru/golang-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
const UserColl = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	coll := client.Database(DBNAME).Collection(UserColl)
	
	return &MongoUserStore{
		client: client,
		coll: coll,
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}