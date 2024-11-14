package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/horlakz/go-auth/payload/response"
)

func GetUserId(c *fiber.Ctx) uuid.UUID {
	userId := c.Locals("userId").(uuid.UUID)

	return userId
}

func Index(c *fiber.Ctx) error {

	var resp response.Response

	var about struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Author  string `json:"author"`
	}

	about.Name = "Go Auth API"
	about.Version = "0.0.1"
	about.Author = "Horlakz"

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"about": about}

	return c.JSON(resp)
}

func NotFound(c *fiber.Ctx) error {
	var resp response.Response

	resp.Status = http.StatusNotFound
	resp.Message = "Route not found"

	return c.Status(http.StatusNotFound).JSON(resp)
}
