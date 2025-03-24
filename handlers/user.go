package handlers

import (
	"errors"
	"log"
	"time"

	"github.com/Guest-01/tryco-backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func convertUserToResponse(u sqlc.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

// @summary		Get User By ID
// @tags		users
// @router		/api/v1/users/{id} [get]
// @param		id	path	int	true	"User ID"
// @success		200	{object}	UserResponse
// @failure		500	{object}	string
// @failure		400	{object}	string
// @failure		404 {object}	string
func (h *Handler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Println("error parsing id: ", err)
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	user, err := h.queries.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(convertUserToResponse(user))
}

// @summary		Create User
// @tags		users
// @router		/api/v1/users [post]
// @param		user body		sqlc.CreateUserParams true	"User object"
// @success		200	{object}	sqlc.CreateUserRow
// @failure		500	{object}	string
// @failure		400	{object}	string
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var user sqlc.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("could not parse body: ", err)
		return c.Status(fiber.StatusBadRequest).SendString("invalid body")
	}

	created, err := h.queries.CreateUser(c.Context(), sqlc.CreateUserParams{
		Email:    user.Email,
		Password: user.Password, // TODO: hash the password
		Username: user.Username,
	})
	if err != nil {
		log.Println("could not create user: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}
