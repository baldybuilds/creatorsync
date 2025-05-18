package server

import (
	"github.com/gofiber/fiber/v2"

	"creatorsync-go/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "creatorsync-go",
			AppName:      "creatorsync-go",
		}),

		db: database.New(),
	}

	return server
}
