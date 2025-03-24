package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/Guest-01/tryco-backend/db/sqlc"
	_ "github.com/Guest-01/tryco-backend/docs"

	"github.com/Guest-01/tryco-backend/handlers"
)

// @title						TryCo API
// @version						1.0
// @description					This is the backend for TryCo.
// @securityDefinitions.apiKey	SessionCookie
// @in							cookie
// @name						_SESSION_ID
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	queries := sqlc.New(conn)

	handler := handlers.New(queries)

	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is the backend for TryCo.")
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// example routes
	v1.Get("/books", handler.GetBooks)
	v1.Get("/books/:id", handler.GetBook)
	v1.Post("/books", handler.CreateBook)
	v1.Put("/books/:id", handler.UpdateBook)
	v1.Delete("/books/:id", handler.DeleteBook)

	log.Fatal(app.Listen(":3000"))
}
