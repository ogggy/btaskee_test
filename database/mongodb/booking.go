package mongodb

import (
	"app/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	BookingAt    time.Time          `json:"bookingAt" bson:"booking_at"`
	CustomerInfo CustomerInfo       `json:"customerInfo" bson:"customer_info"`
	HelperInfo   HelperInfo         `json:"helperInfo" bson:"helper_info"`
	JobInfo      JobInfo            `json:"jobInfo" bson:"job_info"`
}

type CustomerInfo struct {
	CustomerID  primitive.ObjectID `json:"customerID" bson:"customer_id"`
	PhoneNumber string             `json:"phoneNumber" bson:"phone_number"`
	Name        string             `json:"name" bson:"name"`
}

type JobInfo struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

type HelperInfo struct {
	HelperID      primitive.ObjectID `json:"helperID" bson:"helper_id"`
	Name          string             `json:"name" bson:"name"`
	PhoneNumber   string             `json:"phoneNumber" bson:"phone_number"`
	ReceivedJobAt time.Time          `json:"receivedJobAt" bson:"received_job_at"`
	ArrivedAt     time.Time          `json:"arrivedAt" bson:"arrived_at"`
	DoneAt        time.Time          `json:"doneAt" bson:"done_at"`
}

func (ins *bookingDB) New(ctx context.Context, booking *BookingModel) (
	primitive.ObjectID, error) {
	res, err := ins.co.InsertOne(ctx, booking)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (ins *bookingDB) Get(ctx context.Context, id primitive.ObjectID) (
	BookingModel, error) {
	var (
		booking BookingModel
	)
	err := ins.co.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		return booking, err
	}
	return booking, nil
}

func (ins *bookingDB) UpdateHelperReceiveBooking(ctx context.Context,
	bookingID primitive.ObjectID, helperInfo HelperInfo) error {
	var (
		filter = bson.M{
			"_id":                   bookingID,
			"helper_info.helper_id": primitive.NilObjectID,
		}
		update = bson.M{"$set": bson.M{
			"helper_info": helperInfo,
		}}
	)
	_, err := ins.co.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

type bookingDB struct {
	co *mongo.Collection
}

func newBookingDB(db *mongo.Database) *bookingDB {
	return &bookingDB{
		co: database.MongoInit(db, "bookings"),
	}
}
