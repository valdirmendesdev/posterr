package profiles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/domain/users"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/users"
	"github.com/valdirmendesdev/posterr/pkg/services/profiles"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func createUser(t *testing.T) *users.User {
	u, err := users.NewUser(
		types.NewUUID(),
		users.Username("username"),
		time.Now(),
	)
	assert.NoError(t, err)
	return u
}

func TestGetByUsername(t *testing.T) {
	ur := users_infra.NewMemoryRepository()
	u := createUser(t)
	ur.Add(u)

	s := profiles.NewGetByUsernameService(ur)

	p, err := s.Perform(profiles.GetByUsernameRequest{
		Username: users.Username("username"),
	})

	assert.NotNil(t, p)
	assert.NoError(t, err)
}

func TestGetByUsername_user_not_found(t *testing.T) {
	ur := users_infra.NewMemoryRepository()
	s := profiles.NewGetByUsernameService(ur)

	_, err := s.Perform(profiles.GetByUsernameRequest{
		Username: users.Username("notfound"),
	})

	assert.ErrorIs(t, err, users.ErrNotFound)
}
