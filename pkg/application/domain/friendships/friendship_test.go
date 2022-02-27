package friendships_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func TestCreateNewFriendship(t *testing.T) {
	user, err := users.NewUser(types.NewUUID(), users.Username("myuser"), time.Now())
	assert.NoError(t, err)
	follower, err := users.NewUser(types.NewUUID(), users.Username("myfollower"), time.Now())
	assert.NoError(t, err)

	tests := []struct {
		name     string
		user     *users.User
		follower *users.User
		wantErr  bool
		err      error
	}{
		{
			name:     "Valid new friendship",
			user:     user,
			follower: follower,
			wantErr:  false,
			err:      nil,
		},
		{
			name:     "Without user",
			user:     nil,
			follower: nil,
			wantErr:  true,
			err:      friendships.ErrFriendshipInvalidUser,
		},
		{
			name:     "Without follower",
			user:     user,
			follower: nil,
			wantErr:  true,
			err:      friendships.ErrFriendshipInvalidFollower,
		},
		{
			name:     "User trying to follow himself",
			user:     user,
			follower: user,
			wantErr:  true,
			err:      friendships.ErrFriendshipWithYourself,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			friendship, err := friendships.NewFriendship(tt.user, tt.follower)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, friendship.User)
				assert.Equal(t, tt.follower, friendship.Follower)
			}
		})
	}
}
