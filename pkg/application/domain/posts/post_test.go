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

func TestNewPost_invalid_user(t *testing.T) {
	p, err := posts.NewPost(types.NewUUID(), nil, "content", time.Now())
	assert.Nil(t, p)
	assert.Equal(t, posts.ErrInvalidUser, err)
}

func TestNewPost_invalid_created_at(t *testing.T) {
	u := createUser(t)
	var date time.Time
	p, err := posts.NewPost(types.NewUUID(), u, "content", date)
	assert.Nil(t, p)
	assert.Equal(t, posts.ErrInvalidCreatedAt, err)
}

func generateString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

func TestNewPost_content_over_limit(t *testing.T) {
	u := createUser(t)
	content := generateString(posts.ContentMaxLength + 1)
	p, err := posts.NewPost(types.NewUUID(), u, content, time.Now())
	assert.Nil(t, p)
	assert.Equal(t, posts.ErrContentOverMaxLength, err)
}
