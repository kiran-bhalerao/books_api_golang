package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber"
	"github.com/kiranbhalerao123/books_go/book"
	"github.com/kiranbhalerao123/books_go/database"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
}

func init() {
	err := database.Mg.Client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// create app instance
	app := fiber.New()

	// setup the routes
	setupRoutes(app)

	if err := app.Listen(3000); err != nil {
		log.Fatal(err)
	}
}
