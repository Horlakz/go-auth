package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/horlakz/go-auth/constants"
	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/handler"
)

func InitializeRouter(router *fiber.App, dbConn database.DatabaseInterface, env constants.Env) {

	main := router.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	main.Get("/monitor", monitor.New(monitor.Config{Title: "Go Auth API Monitor"}))

	InitializeAuthRouter(main, dbConn, env)

	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	router.Get("/", handler.Index)
	router.Get("*", handler.NotFound)

}
