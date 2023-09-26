package booking

import (
	"app/database/mongodb"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseResp struct {
	ErrCode    int    `json:"errCode"` //errCode == 0 -> success; errCode != 0 error
	ErrMessage string `json:"errMessage"`
}

type CreateBookingReq struct {
	CustomerID primitive.ObjectID `json:"customerID"`
	JobCode    string             `json:"jobCode"`
}

func (req *CreateBookingReq) invalid() error {
	if req.CustomerID.IsZero() {
		return errors.New("customer id invalid")
	}
	if len(req.JobCode) == 0 {
		return errors.New("job code invalid")
	}
	return nil
}

type CreateBookingResp struct {
	BaseResp
	BookingID primitive.ObjectID `json:"bookingID"`
	BookingAt time.Time          `json:"bookingAt"`
}

type GetBookingInfoResp struct {
	BaseResp
	Result mongodb.BookingModel `json:"result"`
}
