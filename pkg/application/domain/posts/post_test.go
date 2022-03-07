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
	assert.Equal(t, posts.TypePost, p.Type)
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

func TestNewRepost(t *testing.T) {
	u := createUser(t)
	now := time.Now()
	p, err := posts.NewPost(types.NewUUID(), u, "content", now)
	assert.NoError(t, err)
	r, err := posts.NewRepost(types.NewUUID(), u, p, now)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.NotNil(t, r.ID)
	assert.Equal(t, u, r.User)
	assert.Equal(t, "", r.Content)
	assert.Equal(t, posts.TypeRepost, r.Type)
	assert.Equal(t, now, r.CreatedAt)
	assert.Equal(t, p, r.ReferencedPost)
}

func TestNewRepost_without_referenced_post(t *testing.T) {
	u := createUser(t)
	r, err := posts.NewRepost(types.NewUUID(), u, nil, time.Now())
	assert.Nil(t, r)
	assert.ErrorIs(t, posts.ErrInvalidReferencedPost, err)
}

func TestNewQuotePost(t *testing.T) {
	u := createUser(t)
	now := time.Now()
	p, err := posts.NewPost(types.NewUUID(), u, "content", now)
	assert.NoError(t, err)
	r, err := posts.NewQuotePost(types.NewUUID(), u, p, "quote", now)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.NotNil(t, r.ID)
	assert.Equal(t, u, r.User)
	assert.Equal(t, "quote", r.Content)
	assert.Equal(t, posts.TypeQuote, r.Type)
	assert.Equal(t, now, r.CreatedAt)
	assert.Equal(t, p, r.ReferencedPost)
}

func TestNewQuotePost_without_comment(t *testing.T) {
	u := createUser(t)
	now := time.Now()
	p, err := posts.NewPost(types.NewUUID(), u, "content", now)
	assert.NoError(t, err)
	r, err := posts.NewQuotePost(types.NewUUID(), u, p, "", now)
	assert.Nil(t, r)
	assert.ErrorIs(t, posts.ErrInvalidContent, err)
}
