package store_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/talvor/asyncapi/dto"
	"github.com/talvor/asyncapi/fixtures"
	"github.com/talvor/asyncapi/store"
)

var _ = Describe("ReportStore", Ordered, func() {
	var env *fixtures.TestEnv
	var reportStore *store.ReportStore
	var userStore *store.UserStore

	BeforeAll(func() {
		te, err := fixtures.NewTestEnv()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(te.ContainerCleanup)
		env = te
		reportStore = store.NewReportStore(env.DB)
		userStore = store.NewUserStore(env.DB)
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

	It("should create a report", func() {
		ctx := context.Background()
		now := time.Now()

		report, err := reportStore.Create(ctx, user.ID, "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(report.UserID).To(Equal(user.ID))
		Expect(report.ReportType).To(Equal("test"))
		Expect(now.UnixNano()).To(BeNumerically("<", report.CreatedAt.UnixNano()))
	})

	It("should update a report", func() {
		ctx := context.Background()
		report, err := reportStore.Create(ctx, user.ID, "test")
		Expect(err).NotTo(HaveOccurred())

		startedAt := report.CreatedAt.Add(time.Minute)
		completedAt := startedAt.Add(time.Minute)
		failedAt := completedAt.Add(time.Minute)
		errorMessage := "test error message"
		downloadURL := "http://example.com"
		downloadURLExpiresAt := failedAt.Add(time.Hour)
		outputFilePath := "/tmp/test"

		report.ReportType = "test2"
		report.StartedAt = &startedAt
		report.CompletedAt = &completedAt
		report.FailedAt = &failedAt
		report.ErrorMessage = &errorMessage
		report.DownloadURL = &downloadURL
		report.DownloadURLExpiresAt = &downloadURLExpiresAt
		report.OutputFilePath = &outputFilePath

		report2, err := reportStore.Update(ctx, report)
		Expect(err).NotTo(HaveOccurred())
		Expect(report2.ReportType).ToNot(Equal("test2"))
		Expect(report2.StartedAt).To(Equal(&startedAt))
		Expect(report2.CompletedAt).To(Equal(&completedAt))
		Expect(report2.FailedAt).To(Equal(&failedAt))
		Expect(report2.ErrorMessage).To(Equal(&errorMessage))
		Expect(report2.DownloadURL).To(Equal(&downloadURL))
		Expect(report2.DownloadURLExpiresAt).To(Equal(&downloadURLExpiresAt))
		Expect(report2.OutputFilePath).To(Equal(&outputFilePath))
	})

	It("should get a report by user id and report id", func() {
		ctx := context.Background()
		report, err := reportStore.Create(ctx, user.ID, "test")
		Expect(err).NotTo(HaveOccurred())

		report2, err := reportStore.ByPrimaryKey(ctx, user.ID, report.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(report2.ID).To(Equal(report.ID))
		Expect(report2.ReportType).To(Equal(report.ReportType))
	})
})
