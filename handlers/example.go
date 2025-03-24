package handlers

import (
	"log"
	"strconv"

	"github.com/Guest-01/tryco-backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

// @summary		GetBooks
// @tags		books
// @router		/api/v1/books [get]
// @success		200	{array}		sqlc.Book
// @failure		500	{object}	fiber.Map
func (h *Handler) GetBooks(c *fiber.Ctx) error {
	books, err := h.queries.GetBooks(c.Context())
	if err != nil {
		log.Println("Error getting books:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get books",
		})
	}

	return c.JSON(books)
}

// @summary		GetBook
// @tags		books
// @router		/api/v1/books/{id} [get]
// @param		id	path	int	true	"Book ID"
// @success	200	{object}	sqlc.Book
// @failure	500	{object}	fiber.Map
// @failure	400	{object}	fiber.Map
func (h *Handler) GetBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Error converting id to int:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}
	book, err := h.queries.GetBook(c.Context(), id)
	if err != nil {
		log.Println("Error getting book:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get book",
		})
	}

	return c.JSON(book)
}

// @summary		CreateBook
// @tags		books
// @router		/api/v1/books [post]
// @success		201	{object}	sqlc.Book
// @failure		500	{object}	fiber.Map
// @failure		400	{object}	fiber.Map
// @param		book	body	sqlc.CreateBookParams	true	"Book object"
// @security    SessionCookie
func (h *Handler) CreateBook(c *fiber.Ctx) error {
	var book sqlc.Book
	if err := c.BodyParser(&book); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	newBook, err := h.queries.CreateBook(c.Context(), sqlc.CreateBookParams{
		Title:  book.Title,
		Author: book.Author,
	})
	if err != nil {
		log.Println("Error creating book:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(newBook)
}

// @summary		UpdateBook
// @tags		books
// @router		/api/v1/books/{id} [put]
// @param		id	path	int	true	"Book ID"
// @param		book	body	sqlc.UpdateBookParams	true	"Book object"
// @success	200	{object}	sqlc.Book
// @failure	500	{object}	fiber.Map
// @failure	400	{object}	fiber.Map
// @security    SessionCookie
func (h *Handler) UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Error converting id to int:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	var book sqlc.Book
	if err := c.BodyParser(&book); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	book.ID = id

	updatedBook, err := h.queries.UpdateBook(c.Context(), sqlc.UpdateBookParams{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
	})
	if err != nil {
		log.Println("Error updating book:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}
	return c.JSON(updatedBook)
}

// @summary		DeleteBook
// @tags		books
// @router		/api/v1/books/{id} [delete]
// @param		id	path	int	true	"Book ID"
// @success		204
// @failure		500	{object}	fiber.Map
// @failure		400	{object}	fiber.Map
// @security    SessionCookie
func (h *Handler) DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Error converting id to int:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	err = h.queries.DeleteBook(c.Context(), id)
	if err != nil {
		log.Println("Error deleting book:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
