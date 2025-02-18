package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talvor/asyncapi/fixtures"
	"github.com/talvor/asyncapi/store"
)

func TestUserStore(t *testing.T) {
	env := fixtures.NewTestEnv(t)
	cleanup := env.SetupDB(t)
	t.Cleanup(func() {
		cleanup(t)
	})

	ctx := context.Background()
	now := time.Now()

	userStore := store.NewUserStore(env.DB)
	user, err := userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
	require.NoError(t, err)
	require.Less(t, now.UnixNano(), user.CreatedAt.UnixNano())

	require.Equal(t, user.Email, "test@testing.com")
	require.NoError(t, user.ComparePassword("testingpassword"))

	user2, err := userStore.ByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
}
