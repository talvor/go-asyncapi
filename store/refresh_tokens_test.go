package store_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/talvor/asyncapi/apiserver"
	"github.com/talvor/asyncapi/dto"
	"github.com/talvor/asyncapi/fixtures"
	"github.com/talvor/asyncapi/store"
)

var _ = Describe("RefreshTokenStore", Ordered, func() {
	var env *fixtures.TestEnv
	var refreshTokenStore *store.RefreshTokenStore
	var userStore *store.UserStore
	var jwtManager *apiserver.JwtManager

	BeforeAll(func() {
		te, err := fixtures.NewTestEnv()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(te.ContainerCleanup)
		env = te
		refreshTokenStore = store.NewRefreshTokenStore(env.DB)
		userStore = store.NewUserStore(env.DB)
		jwtManager = apiserver.NewJwtManager(env.Config)
	})

	var user *dto.User
	BeforeEach(func() {
		err := env.SetupDB()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(env.TeardownDB)

		ctx := context.Background()
		user, err = userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should create a refresh token", func() {
		ctx := context.Background()
		now := time.Now()

		tokenPair, err := jwtManager.GenerateTokenPair(user.ID)
		Expect(err).NotTo(HaveOccurred())

		expiresAt, err := tokenPair.RefreshToken.Claims.GetExpirationTime()
		Expect(err).NotTo(HaveOccurred())

		refreshToken, err := refreshTokenStore.Create(ctx, user.ID, tokenPair.RefreshToken)
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshToken.UserID).To(Equal(user.ID))
		Expect(refreshToken.HashedToken).NotTo(BeEmpty())
		Expect(now.UnixNano()).To(BeNumerically("<", refreshToken.CreatedAt.UnixNano()))
		Expect(refreshToken.ExpiresAt.UnixMilli()).To(Equal(expiresAt.UnixMilli()))
	})

	It("should retrieve a refresh token by user id and token", func() {
		ctx := context.Background()

		tokenPair, err := jwtManager.GenerateTokenPair(user.ID)
		Expect(err).NotTo(HaveOccurred())

		refreshToken1, err := refreshTokenStore.Create(ctx, user.ID, tokenPair.RefreshToken)
		Expect(err).NotTo(HaveOccurred())

		refreshToken2, err := refreshTokenStore.ByPrimaryKey(ctx, user.ID, tokenPair.RefreshToken)
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshToken2.UserID).To(Equal(refreshToken1.UserID))
		Expect(refreshToken2.HashedToken).To(Equal(refreshToken1.HashedToken))
		Expect(refreshToken2.CreatedAt).To(Equal(refreshToken1.CreatedAt))
		Expect(refreshToken2.ExpiresAt).To(Equal(refreshToken1.ExpiresAt))

	})
	It("should delete all user refresh tokens", func() {
		ctx := context.Background()

		tokenPair, err := jwtManager.GenerateTokenPair(user.ID)
		Expect(err).NotTo(HaveOccurred())

		_, err = refreshTokenStore.Create(ctx, user.ID, tokenPair.RefreshToken)
		Expect(err).NotTo(HaveOccurred())

		result, err := refreshTokenStore.DeleteUserTokens(ctx, user.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RowsAffected()).To(Equal(int64(1)))
	})
})
