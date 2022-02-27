package utils

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func GetLoggedUser() *users.User {
	id, _ := types.ParseUUID("d488526d-a8df-43e6-a226-a2d355adc0e5")
	return &users.User{
		ID:       id,
		Username: "loggedUser",
		JoinedAt: time.Now(),
	}
}
