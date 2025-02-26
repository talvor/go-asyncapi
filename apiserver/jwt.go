package apiserver

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/talvor/asyncapi/config"
)

var signingMethod = jwt.SigningMethodHS256

type JwtManager struct {
	config *config.Config
}

func NewJwtManager(config *config.Config) *JwtManager {
	return &JwtManager{
		config: config,
	}
}

type TokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (j *JwtManager) Parse(token string) (*jwt.Token, error) {
	parser := jwt.NewParser()

	jwtToken, err := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.config.JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return jwtToken, nil
}

func (j *JwtManager) IsAccessToken(token *jwt.Token) bool {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	tokenType, ok := claims["token_type"].(string)
	if !ok {
		return false
	}
	return tokenType == "access"
}

func (j *JwtManager) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	now := time.Now()
	issuer := "http://" + j.config.APIHost + ":" + j.config.APIPort
	claims := CustomClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Generate access token
	accessToken, err := j.GenerateToken(&claims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	claims.TokenType = "refresh"
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Hour * 24 * 30))
	refreshToken, err := j.GenerateToken(&claims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JwtManager) GenerateToken(claims *CustomClaims) (*jwt.Token, error) {
	jwtToken := jwt.NewWithClaims(signingMethod, claims)

	key := []byte(j.config.JwtSecret)
	var err error
	jwtToken.Raw, err = jwtToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign %s token: %w", claims.TokenType, err)
	}
	return jwtToken, nil
}

func (j *JwtManager) GetUserIDFromToken(token *jwt.Token) (uuid.UUID, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.New(), fmt.Errorf("failed to get subject from token: %w", err)
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.New(), fmt.Errorf("failed to convert subject to UUID: %w", err)
	}

	return userID, nil
}
