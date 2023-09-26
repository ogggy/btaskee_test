package main

import (
	"app/api/booking"
	"app/api/pricing"
	"app/api/sending"
	"app/database/mongodb"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ReadTimeoutSecond  = 30
	WriteTimeoutSecond = 30
	IdleTimeoutSecond  = 60
	DbURI              = "mongodb://localhost:27017"
	DbName             = "btaskee"
	DbTimeout          = 15 * time.Second
)

func listenSignal() chan os.Signal {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	return s
}

func startAPIServer(interrupt chan<- os.Signal, bindAddress string) {
	// Make new server engine
	apiEngine := func() *gin.Engine {
		engine := gin.New()
		engine.Use(
			gin.LoggerWithWriter(os.Stdout),
		)
		return engine
	}()

	dbStore := mongodb.NewStoreDB(DbURI, DbName, DbTimeout)
	//add example data
	go func() {
		ctx := context.Background()
		if err := dbStore.JobDB.InsertFromFile(ctx,
			"./example_data/jobs.json"); err != nil {
			log.Println("[InsertFromFile-MongoDB-Err]", err.Error())
		}
		if err := dbStore.CustomerDB.InsertFromFile(ctx,
			"./example_data/customers.json"); err != nil {
			log.Println("[InsertFromFile-MongoDB-Err]", err.Error())
		}
		if err := dbStore.HelperDB.InsertFromFile(ctx,
			"./example_data/helpers.json"); err != nil {
			log.Println("[InsertFromFile-MongoDB-Err]", err.Error())
		}
	}()

	//create service
	pricingService := pricing.NewService()
	sendingService := sending.NewService(dbStore)
	bookingService := booking.NewService(dbStore, sendingService)

	//apply api routes
	pricing.NewHandler(pricingService).Apply(apiEngine)
	booking.NewHandler(bookingService).Apply(apiEngine)

	go func() {
		server := http.Server{
			Addr:         bindAddress,
			Handler:      apiEngine,
			ReadTimeout:  ReadTimeoutSecond * time.Second,
			WriteTimeout: WriteTimeoutSecond * time.Second,
			IdleTimeout:  IdleTimeoutSecond * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			println(
				"\r\n"+
					"ListenAndServe:", err.Error())
		}
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			println(
				"\r\n"+
					"Server Shutdown:", err.Error(),
			)
		}
		interrupt <- os.Interrupt
	}()
	println(
		"[ServerAPIs] Listening:", bindAddress,
	)
}
