package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetTwitchStreamsHandler(c *fiber.Ctx) error {
	// TO DO: implement getTwitchStreamsHandler
	return c.JSON(fiber.Map{
		"message": "getTwitchStreamsHandler not implemented",
	})
}
