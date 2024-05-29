package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mayank12gt/free-audiobooks-backend/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type app struct {
	logger   *log.Logger
	services services.Services
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if $PORT is not set
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if err := godotenv.Load(); err != nil {
		log.Print("no env file found")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Print("No DSN found")
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Client().Disconnect(context.TODO()); err != nil {
			log.Panic(err)
		}
	}()

	app := &app{
		logger:   logger,
		services: services.NewService(db),
	}

	err = app.serve(port)
	if err != nil {
		app.logger.Fatal(err)
	}

}

func (app *app) serve(port string) error {
	server := echo.New()
	//server.Use(middleware.CORS())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{

		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	app.registerHandlers(server)
	err := server.Start(":" + port)
	app.logger.Printf("server started")

	if err != nil {
		return err
	}

	return nil

}

func (app *app) registerHandlers(server *echo.Echo) {
	server.GET("/audiobooks", app.listHandler())

	server.GET("/audiobooks/:id", app.GetHandler())

}

func openDB(dsn string) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	log.Print("DB connected")
	return client.Database("audiobooksDB"), nil
}
