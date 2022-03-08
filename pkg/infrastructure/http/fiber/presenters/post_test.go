package presenters_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

func TestTransformPostToPresenter(t *testing.T) {
	u := utils.GetLoggedUser()
	post, err := posts.NewPost(types.NewUUID(), u, "my post", time.Now())
	assert.NoError(t, err)
	presenter, err := presenters.TransformPostToPresenter(post)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, presenter.ID)
	assert.Equal(t, post.User.Username, presenter.Username)
	assert.Equal(t, post.Content, presenter.Content)
	assert.Equal(t, post.Type, presenter.Type)
	assert.Equal(t, post.CreatedAt, presenter.CreatedAt)
}
