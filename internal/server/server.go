package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"creatorsync-go/internal/clerk"
	"creatorsync-go/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() (*FiberServer, error) {
	if err := clerk.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize Clerk client: %w", err)
	}

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "creatorsync",
			AppName:      "creatorsync",
		}),

		db: database.New(),
	}

	return server, nil
}
