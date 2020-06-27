package book

import (
	"context"

	"github.com/gofiber/fiber"
	"github.com/kiranbhalerao123/books_go/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(c *fiber.Ctx) {
	db := database.Mg.Db

	var books []Book = make([]Book, 0)

	// get all records as a cursor
	query := bson.D{{}}
	cursor, err := db.Collection("books").Find(c.Fasthttp, query)

	// pipe := []bson.M{
	// 	{
	// 		"$match": bson.M{
	// 			"author": "kiran",
	// 		},
	// 	},
	// }
	// cursor, err := db.Collection("books").Aggregate(c.Fasthttp, pipe)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// iterate the cursor and decode each item into books
	if err := cursor.All(c.Fasthttp, &books); err != nil {
		c.Status(500).Send(err)
		return
	}

	//dont forget to close the cursor
	defer cursor.Close(context.Background())

	if err := c.JSON(books); err != nil {
		c.Status(500).Send(err)
		return
	}
}

func GetBook(c *fiber.Ctx) {
	db := database.Mg.Db

	id := c.Params("id")
	bookId, err := primitive.ObjectIDFromHex(id)

	// the provided ID might be invalid ObjectID
	if err != nil {
		c.Status(400).Send()
		return
	}

	// Find the employee and update its data
	query := bson.D{{Key: "_id", Value: bookId}}

	var book Book

	err = db.Collection("books").FindOne(c.Fasthttp, query).Decode(&book)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// return the Book in JSON format
	if err := c.Status(201).JSON(book); err != nil {
		c.Status(500).Send(err)
		return
	}
}

func NewBook(c *fiber.Ctx) {
	db := database.Mg.Db

	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		c.Status(503).Send(err)
		return
	}

	// force MongoDB to always set its own generated ObjectIDs
	book.ID = ""

	// insert the record
	insertionResult, err := db.Collection("books").InsertOne(c.Fasthttp, book)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// get the just inserted record in order to return it as response
	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}

	var createdBook Book
	err = db.Collection("books").FindOne(c.Fasthttp, filter).Decode(&createdBook)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// return the created Book in JSON format
	if err := c.Status(201).JSON(createdBook); err != nil {
		c.Status(500).Send(err)
		return
	}
}

func DeleteBook(c *fiber.Ctx) {
	db := database.Mg.Db

	id := c.Params("id")
	bookId, err := primitive.ObjectIDFromHex(id)

	// the provided ID might be invalid ObjectID
	if err != nil {
		c.Status(400).Send()
		return
	}

	// find and delete the Book with the given ID
	query := bson.D{{Key: "_id", Value: bookId}}
	result, err := db.Collection("books").DeleteOne(c.Fasthttp, &query)

	if err != nil {
		c.Status(500).Send()
		return
	}

	// the Book might not exist
	if result.DeletedCount < 1 {
		c.Status(404).Send()
		return
	}

	c.Status(204).Send("Book Successfully deleted")
}
