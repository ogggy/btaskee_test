package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(uri, dbname string, timeout time.Duration) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbname)
	if err := client.Ping(ctx, db.ReadPreference()); err != nil {
		return nil, err
	}
	return db, nil
}

/*
MongoIndex object
*/
type MongoIndex struct {
	Keys   interface{}
	Unique bool
}

/*
MongoInit : create first collection config
*/
func MongoInit(db *mongo.Database, collectionname string, index ...MongoIndex) *mongo.Collection {
	var (
		collection *mongo.Collection
		indexes    = []mongo.IndexModel{}
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	var collectionValidate = func() (created bool) {
		list, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			log.Println(err)
			return
		}
		for _, n := range list {
			if n == collectionname {
				created = true
				break
			}
		}
		return
	}
	if created := collectionValidate(); created {
		log.Printf("Log-Debug: Collection `%s` is already available. collectionValidate=%v\r\n", collectionname, created)
	} else {
		if err := db.CreateCollection(ctx, collectionname); err != nil {
			log.Printf("Log-Error: %+v\r\n", err)
		} else {
			log.Printf("Log-Debug: Collection `%s` is already available\r\n", collectionname)
		}

	}
	collection = db.Collection(collectionname)

	for _, uq := range index {
		if uq.Keys == nil {
			continue
		}
		indexes = append(indexes, mongo.IndexModel{
			Keys:    uq.Keys,
			Options: options.Index().SetUnique(uq.Unique),
		})
	}

	if len(indexes) > 0 {
		names, err := collection.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			log.Printf("Log-Error: %+v\r\n", err)

		}
		for _, name := range names {
			log.Printf("Log-Debug: Index created `%s`\r\n", name)
		}
	}

	return collection
}
