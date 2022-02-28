package posts_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	posts_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func createUser(t *testing.T, username string) *users.User {
	u, err := users.NewUser(types.NewUUID(), users.Username(username), time.Now())
	assert.NoError(t, err)
	return u
}

func createPost(t *testing.T, u *users.User, createdAt time.Time) *posts.Post {
	p, err := posts.NewPost(types.NewUUID(), u, "my post", createdAt)
	assert.NoError(t, err)
	return p
}

func TestInsertPost(t *testing.T) {
	u := createUser(t, "myuser")
	p := createPost(t, u, time.Now())
	r := posts_infra.NewMemoryRepository()
	err := r.Insert(p)
	assert.NoError(t, err)

	posts, err := r.List()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
}

func TestGetPostsByUsername(t *testing.T) {
	r := posts_infra.NewMemoryRepository()
	u := createUser(t, "myuser")
	p := createPost(t, u, time.Now())
	err := r.Insert(p)
	assert.NoError(t, err)

	u = createUser(t, "otheruser")
	p = createPost(t, u, time.Now())
	err = r.Insert(p)
	assert.NoError(t, err)

	posts, err := r.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	posts, err = r.ListByUsername("myuser")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
}

func TestGetPostsByDate(t *testing.T) {
	createdAt := time.Now()
	r := posts_infra.NewMemoryRepository()
	u := createUser(t, "myuser")
	p := createPost(t, u, createdAt)
	err := r.Insert(p)
	assert.NoError(t, err)
	yesterday := createdAt.AddDate(0, 0, -1)
	p = createPost(t, u, yesterday)
	err = r.Insert(p)
	assert.NoError(t, err)

	posts, err := r.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	posts, err = r.ListByDate(createdAt)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
}
