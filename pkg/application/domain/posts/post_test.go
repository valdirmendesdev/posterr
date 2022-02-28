package posts_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func createUser(t *testing.T) *users.User {
	u, err := users.NewUser(types.NewUUID(), users.Username("username"), time.Now())
	assert.NoError(t, err)
	return u
}

func TestNewPost(t *testing.T) {
	u := createUser(t)
	now := time.Now()
	p, err := posts.NewPost(types.NewUUID(), u, "content", now)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.NotNil(t, p.ID)
	assert.Equal(t, u, p.User)
	assert.Equal(t, "content", p.Content)
	assert.Equal(t, now, p.CreatedAt)
}
