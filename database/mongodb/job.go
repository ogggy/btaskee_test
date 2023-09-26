package mongodb

import (
	"app/database"
	"context"
	"encoding/json"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobModel struct {
	Code string `json:"Code" bson:"code"`
	Name string `json:"Name" bson:"name"`
}

func (ins *jobDB) InsertFromFile(ctx context.Context, path string) error {

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

func (ins *jobDB) GetByCode(ctx context.Context, jobCode string) (
	JobModel, error) {
	var (
		job    JobModel
		filter = bson.M{"code": jobCode}
	)
	if err := ins.co.FindOne(ctx, filter).Decode(&job); err != nil {
		return JobModel{}, err
	}
	return job, nil
}

type jobDB struct {
	co *mongo.Collection
}

func newJobDB(db *mongo.Database) *jobDB {
	return &jobDB{
		co: database.MongoInit(db, "jobs",
			database.MongoIndex{
				Keys:   "code",
				Unique: true,
			}),
	}
}
