package booking

import (
	"app/api/sending"
	"app/database/mongodb"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	sendingService *sending.Service
	db             *mongodb.Store
}

func (s *Service) CreateNewBooking(ctx context.Context, req *CreateBookingReq) (
	CreateBookingResp, error) {

	var (
		resp CreateBookingResp
		// err  error
	)

	if err := req.invalid(); err != nil {
		resp.ErrCode = 40
		resp.ErrMessage = err.Error()
		return resp, err
	}

	// validate job
	job, err := s.db.JobDB.GetByCode(ctx, req.JobCode)
	if errors.Is(err, mongo.ErrNoDocuments) {
		resp.ErrCode = 44
		resp.ErrMessage = "no job match"
		return resp, err
	}
	if err != nil {
		resp.ErrCode = 53
		resp.ErrMessage = "database error"
		return resp, err
	}
	//validate customer
	cus, err := s.db.CustomerDB.Get(ctx, req.CustomerID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		resp.ErrCode = 44
		resp.ErrMessage = "not found customer"
		return resp, err
	}
	if err != nil {
		resp.ErrCode = 53
		resp.ErrMessage = "database error"
		return resp, err
	}

	//save booking
	newBooking := &mongodb.BookingModel{
		ID:        primitive.NewObjectID(),
		BookingAt: time.Now(),
		CustomerInfo: mongodb.CustomerInfo{
			CustomerID:  req.CustomerID,
			PhoneNumber: cus.PhoneNumber,
			Name:        cus.Name,
		},
		JobInfo: mongodb.JobInfo{
			Code: job.Code,
			Name: job.Name,
		},
	}
	bookingID, err := s.db.BookingDB.New(ctx, newBooking)
	if err != nil {
		resp.ErrCode = 53
		resp.ErrMessage = "database error"
		return resp, err
	}
	resp.BookingID = bookingID
	resp.BookingAt = newBooking.BookingAt

	//TODO: send to helper
	go s.sendingService.SendJobToHelper(context.Background(), bookingID)

	return resp, nil
}

func (s *Service) GetBookingInfo(ctx context.Context, id primitive.ObjectID) (
	GetBookingInfoResp, error) {
	var (
		resp GetBookingInfoResp
	)

	bookingInf, err := s.db.BookingDB.Get(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		resp.ErrCode = 44
		resp.ErrMessage = "not found booking"
		return resp, err
	}
	if err != nil {
		resp.ErrCode = 53
		resp.ErrMessage = "database error"
		return resp, err
	}
	resp.Result = bookingInf

	return resp, nil
}

func NewService(db *mongodb.Store, sendingService *sending.Service) *Service {
	return &Service{
		db:             db,
		sendingService: sendingService,
	}
}
