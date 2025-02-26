package apiserver_test

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/talvor/asyncapi/apiserver"
	"github.com/talvor/asyncapi/config"
)

var _ = Describe("JwtManager", Ordered, func() {
	var jwtManager *apiserver.JwtManager

	BeforeAll(func() {
		conf := config.GetConfig()
		jwtManager = apiserver.NewJwtManager(conf)
	})

	It("should generate token pair", func() {
		userID := uuid.New()
		tokenPair, err := jwtManager.GenerateTokenPair(userID)
		Expect(err).NotTo(HaveOccurred())

		Expect(tokenPair.AccessToken).NotTo(BeNil())
		Expect(tokenPair.RefreshToken).NotTo(BeNil())
	})

	It("should parse token", func() {
		userID := uuid.New()
		tokenPair, err := jwtManager.GenerateTokenPair(userID)
		Expect(err).NotTo(HaveOccurred())

		accessToken, err := jwtManager.Parse(tokenPair.AccessToken.Raw)
		Expect(err).NotTo(HaveOccurred())
		Expect(jwtManager.IsAccessToken(accessToken)).To(BeTrue())

		refreshToken, err := jwtManager.Parse(tokenPair.RefreshToken.Raw)
		Expect(err).NotTo(HaveOccurred())
		Expect(jwtManager.IsAccessToken(refreshToken)).To(BeFalse())
	})

	It("should create token for user", func() {
		userID := uuid.New()
		tokenPair, err := jwtManager.GenerateTokenPair(userID)
		Expect(err).NotTo(HaveOccurred())

		subject, err := tokenPair.AccessToken.Claims.GetSubject()
		Expect(err).NotTo(HaveOccurred())
		Expect(subject).To(Equal(userID.String()))

		subject, err = tokenPair.RefreshToken.Claims.GetSubject()
		Expect(err).NotTo(HaveOccurred())
		Expect(subject).To(Equal(userID.String()))
	})

	It("should fail to parse expired token", func() {
		now := time.Now()

		claims := apiserver.CustomClaims{
			TokenType: "access",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "subject",
				Issuer:    "issuer",
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * -15)),
				IssuedAt:  jwt.NewNumericDate(now.Add(time.Minute * -30)),
			},
		}

		token, err := jwtManager.GenerateToken(&claims)
		Expect(err).NotTo(HaveOccurred())

		_, err = jwtManager.Parse(token.Raw)
		Expect(err).Should(MatchError(ContainSubstring("token is expired")))
	})
})
