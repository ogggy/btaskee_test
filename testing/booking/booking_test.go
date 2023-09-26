package booking_test

import (
	"app/api/booking"
	"app/api/sending"
	"app/database/mongodb"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	dbStore        *mongodb.Store
	sendingService *sending.Service
	bookingService *booking.Service
)

func init() {
	dbStore = mongodb.NewStoreDB(
		"mongodb://localhost:27017",
		"btaskee",
		15*time.Second,
	)
	//create service
	sendingService = sending.NewService(dbStore)
	bookingService = booking.NewService(dbStore, sendingService)
}

func TestBooking(t *testing.T) {
	type testcase struct {
		name string
		req  booking.CreateBookingReq
		resp booking.CreateBookingResp
		// err error
	}

	oid, _ := primitive.ObjectIDFromHex("6511fcc730966e90c30740d4")

	tt := []testcase{
		{
			name: "Happy",
			req: booking.CreateBookingReq{
				CustomerID: oid,
				JobCode:    "001",
			},
			resp: booking.CreateBookingResp{
				BaseResp: booking.BaseResp{
					ErrCode: 0,
				},
			},
		},
		{
			name: "Not found customer",
			req: booking.CreateBookingReq{
				CustomerID: primitive.NewObjectID(),
				JobCode:    "001",
			},
			resp: booking.CreateBookingResp{
				BaseResp: booking.BaseResp{
					ErrCode: 44,
				},
			},
		},
		{
			name: "Not found job",
			req: booking.CreateBookingReq{
				CustomerID: oid,
				JobCode:    "003",
			},
			resp: booking.CreateBookingResp{
				BaseResp: booking.BaseResp{
					ErrCode: 44,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			actualResp, err := bookingService.CreateNewBooking(ctx, &tc.req)

			if tc.name != "Happy" {
				assert.Error(t, err)
				assert.NotEqual(t, actualResp.ErrCode, 0)
				assert.Equal(t, actualResp.ErrCode, tc.resp.ErrCode)
				assert.True(t, actualResp.BookingID.IsZero())
				assert.True(t, actualResp.BookingAt.IsZero())
			} else {
				//
				assert.NoError(t, err)
				assert.Equal(t, actualResp.ErrCode, 0)
				assert.False(t, actualResp.BookingID.IsZero())
				assert.False(t, actualResp.BookingAt.IsZero())

				//
				time.Sleep(5 * time.Second)
				//check send data
				bookingInf, _ := dbStore.BookingDB.Get(ctx, actualResp.BookingID)
				assert.Equal(t, actualResp.BookingID, bookingInf.ID)
				assert.False(t, bookingInf.BookingAt.IsZero())
				assert.True(t, bookingInf.BookingAt.Before(time.Now()))
				assert.False(t, bookingInf.CustomerInfo.CustomerID.IsZero())
				assert.NotEmpty(t, bookingInf.JobInfo.Code)
				assert.False(t, bookingInf.HelperInfo.HelperID.IsZero())
				assert.False(t, bookingInf.HelperInfo.ReceivedJobAt.IsZero())

				customer, _ := dbStore.CustomerDB.Get(ctx, bookingInf.CustomerInfo.CustomerID)
				assert.Equal(t, bookingInf.CustomerInfo.CustomerID, customer.ID)

				helper, _ := dbStore.HelperDB.Get(ctx, bookingInf.HelperInfo.HelperID)
				assert.Equal(t, bookingInf.HelperInfo.HelperID, helper.ID)
			}
		})
	}
}
