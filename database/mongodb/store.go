package mongodb

import (
	"app/database"
	"time"
)

type Store struct {
	BookingDB  *bookingDB
	HelperDB   *helperDB
	JobDB      *jobDB
	CustomerDB *customerDB
}

func NewStoreDB(uri, dbname string, timeout time.Duration) *Store {

	db, err := database.MongoConnect(uri, dbname, timeout)
	if err != nil {
		println("\r[MONGO-DB-] MongoConnect to :", uri, ", dbname :", dbname, ", timeout :", timeout,
			"err :", err.Error(),
		)
		panic(err)
	}

	return &Store{
		BookingDB:  newBookingDB(db),
		HelperDB:   newHelperDB(db),
		JobDB:      newJobDB(db),
		CustomerDB: newCustomerDB(db),
	}
}
