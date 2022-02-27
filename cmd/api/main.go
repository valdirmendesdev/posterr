package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/handlers"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func main() {

	app := fiber.New()
	userRepo := users_infra.NewMemoryRepository()
	user, _ := users.NewUser(types.NewUUID(), users.Username("valdirmendes"), time.Now())
	userRepo.Add(user)
	profileConfig := handlers.NewProfileRoutesConfigs(app, userRepo)

	handlers.MountProfilesRoutes(profileConfig)
	app.Listen(":3000")
}
