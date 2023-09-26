package sending

import (
	"app/database/mongodb"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	db *mongodb.Store
}

func (s *Service) SendJobToHelper(ctx context.Context, bookingID primitive.ObjectID) error {
	// Notes: My idea is sending to multiple helpers matched job
	// 		  But I have limitations of time and knowledges
	// ---->  SO I choose random helper to receive job

	bookingInf, err := s.db.BookingDB.Get(ctx, bookingID)
	if err != nil {
		log.Println("[SendJobToHelper-Service-Error]", err)
		return err
	}

	helper, err := s.db.HelperDB.FindHelperMatchJob(ctx, bookingInf.JobInfo.Code)
	if err != nil {
		log.Println("[SendJobToHelper-Service-Error]", err)
		return err
	}

	err = s.db.BookingDB.UpdateHelperReceiveBooking(ctx, bookingID,
		mongodb.HelperInfo{
			HelperID:      helper.ID,
			Name:          helper.Name,
			PhoneNumber:   helper.PhoneNumber,
			ReceivedJobAt: time.Now(),
		})
	if err != nil {
		log.Println("[SendJobToHelper-Service-Error]", err)
		return err
	}

	return nil
}

func NewService(db *mongodb.Store) *Service {
	return &Service{db: db}
}
