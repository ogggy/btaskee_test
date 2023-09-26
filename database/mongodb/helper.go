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

type HelperModel struct {
	ID          primitive.ObjectID `json:"-" bson:"_id"`
	Jobs        []string           `json:"jobs" bson:"jobs"` //list of job_code helper can work
	Name        string             `json:"name" bson:"name"`
	PhoneNumber string             `json:"phoneNumber" bson:"phone_number"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
}

func (ins *helperDB) InsertFromFile(ctx context.Context, path string) error {

	buff, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	data := []interface{}{}
	if err := json.Unmarshal(buff, &data); err != nil {
		return err
	}

	_, err = ins.co.InsertMany(ctx, data)
	if err != nil {
		return err
	}

	return nil

}

func (ins *helperDB) Get(ctx context.Context, id primitive.ObjectID) (
	HelperModel, error) {
	var helper HelperModel

	err := ins.co.FindOne(ctx, bson.M{"_id": id}).Decode(&helper)
	if err != nil {
		return helper, err
	}

	return helper, nil
}

func (ins *helperDB) FindHelperMatchJob(ctx context.Context, jobCode string) (
	HelperModel, error) {
	var (
		helper HelperModel
		filter = bson.M{
			"jobs": jobCode,
		}
	)
	if err := ins.co.FindOne(ctx, filter).Decode(&helper); err != nil {
		return helper, err
	}
	return helper, nil
}

type helperDB struct {
	co *mongo.Collection
}

func newHelperDB(db *mongo.Database) *helperDB {
	return &helperDB{
		co: database.MongoInit(db, "helpers",
			database.MongoIndex{
				Keys:   "phone_number",
				Unique: true,
			}),
	}
}
