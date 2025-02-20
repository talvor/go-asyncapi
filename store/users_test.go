package store_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/talvor/asyncapi/fixtures"
	"github.com/talvor/asyncapi/store"
)

var _ = Describe("UserStore", Ordered, func() {
	var env *fixtures.TestEnv
	var userStore *store.UserStore

	BeforeAll(func() {
		te, err := fixtures.NewTestEnv()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(te.ContainerCleanup)
		env = te
		userStore = store.NewUserStore(env.DB)
	})

	BeforeEach(func() {
		err := env.SetupDB()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(env.TeardownDB)
	})

	It("should create a user", func() {
		ctx := context.Background()
		now := time.Now()

		user, err := userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
		Expect(err).NotTo(HaveOccurred())
		Expect(user.ID).NotTo(BeZero())
		Expect(user.Email).To(Equal("test@testing.com"))
		Expect(now.UnixNano()).To(BeNumerically("<", user.CreatedAt.UnixNano()))
	})

	It("should retrieve a user by ID", func() {
		ctx := context.Background()

		user1, err := userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
		Expect(err).NotTo(HaveOccurred())

		user2, err := userStore.ByID(ctx, user1.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(user2.ID).To(Equal(user1.ID))
		Expect(user2.Email).To(Equal(user1.Email))
		Expect(user2.HashedPasswordBase64).To(Equal(user1.HashedPasswordBase64))
	})

	It("should retrieve a user by Email", func() {
		ctx := context.Background()

		user1, err := userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
		Expect(err).NotTo(HaveOccurred())

		user2, err := userStore.ByEmail(ctx, user1.Email)
		Expect(err).NotTo(HaveOccurred())
		Expect(user2.ID).To(Equal(user1.ID))
		Expect(user2.Email).To(Equal(user1.Email))
		Expect(user2.HashedPasswordBase64).To(Equal(user1.HashedPasswordBase64))
	})

	Context("when a user exists", func() {
		It("should compare correct password", func() {
			ctx := context.Background()

			user, err := userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
			Expect(err).NotTo(HaveOccurred())

			err = user.ComparePassword("testingpassword")
			Expect(err).NotTo(HaveOccurred())
		})
		It("should not compare incorrect password", func() {
			ctx := context.Background()

			user, err := userStore.CreateUser(ctx, "test@testing.com", "testingpassword")
			Expect(err).NotTo(HaveOccurred())

			err = user.ComparePassword("incorrectpassword")
			Expect(err).To(HaveOccurred())
		})
	})
})
