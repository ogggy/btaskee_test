package mongodb

import (
	"app/database"
	"context"
	"encoding/json"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	PhoneNumber string             `json:"phoneNumber" bson:"phone_number"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
}

func (ins *customerDB) InsertFromFile(ctx context.Context, path string) error {

	buff, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cus CustomerModel
	if err := json.Unmarshal(buff, &cus); err != nil {
		return err
	}

	_, err = ins.co.InsertOne(ctx, cus)
	if err != nil {
		return err
	}

	return nil

}

func (ins *customerDB) Get(ctx context.Context, id primitive.ObjectID) (
	CustomerModel, error) {
	var (
		cus    CustomerModel
		filter = bson.M{"_id": id}
	)
	if err := ins.co.FindOne(ctx, filter).Decode(&cus); err != nil {
		return CustomerModel{}, err
	}
	return cus, nil
}

type customerDB struct {
	co *mongo.Collection
}

func newCustomerDB(db *mongo.Database) *customerDB {
	return &customerDB{
		co: database.MongoInit(db, "customers",
			database.MongoIndex{
				Keys:   "phone_number",
				Unique: true,
			}),
	}
}
