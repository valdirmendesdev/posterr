package users_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	users_domain "github.com/valdirmendesdev/posterr/pkg/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func TestGetByUsername(t *testing.T) {
	r := users.NewMemoryRepository()
	u, err := users_domain.NewUser(types.NewUUID(), "validuser", time.Now())
	assert.NoError(t, err)
	err = r.Add(u)
	assert.NoError(t, err)

	user, err := r.GetByUsername(users_domain.Username("validuser"))
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestGetByUsername_NotFound(t *testing.T) {
	r := users.NewMemoryRepository()
	_, err := r.GetByUsername(users_domain.Username("notfounduser"))
	assert.ErrorIs(t, err, users_domain.ErrNotFound)
}
