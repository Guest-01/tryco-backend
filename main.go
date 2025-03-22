package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/Guest-01/tryco-backend/docs"

	"github.com/Guest-01/tryco-backend/handlers"
)

// @title TryCo API
// @version 1.0
// @description This is the backend for TryCo.
func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is the backend for TryCo.")
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/hello", handlers.Hello)

	log.Fatal(app.Listen(":3000"))
}
