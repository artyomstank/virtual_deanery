// pkg/jwt/jwt.go
package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT token claims.
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwtlib.RegisteredClaims
}

// TokenClient provides JWT token operations.
type TokenClient interface {
	// GenerateAccessToken creates new access token.
	GenerateAccessToken(userID int64, email string) (string, error)

	// GenerateRefreshToken creates new refresh token.
	GenerateRefreshToken(userID int64, email string) (string, error)

	// ValidateToken validates token and returns claims.
	ValidateToken(tokenString string) (*Claims, error)

	// ExtractClaims extracts claims from valid token.
	ExtractClaims(tokenString string) (*Claims, error)
}

// JWTClient implements TokenClient.
type JWTClient struct {
	secretKey  string
	algo       string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewJWTClient creates new JWT client.
func NewJWTClient(
	secretKey string,
	algo string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *JWTClient {
	return &JWTClient{
		secretKey:  secretKey,
		algo:       algo,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// GenerateAccessToken creates new access token.
func (jc *JWTClient) GenerateAccessToken(userID int64, email string) (string, error) {
	// TODO: Create claims with AccessTTL expiration
	// claims.Type = "access"

	// TODO: Create token with claims

	// TODO: Sign token using jc.algo and jc.secretKey

	// TODO: Return signed token string or error
	return "", fmt.Errorf("not implemented")
}

// GenerateRefreshToken creates new refresh token.
func (jc *JWTClient) GenerateRefreshToken(userID int64, email string) (string, error) {
	// TODO: Create claims with RefreshTTL expiration
	// claims.Type = "refresh"

	// TODO: Create and sign token

	// TODO: Return signed token string or error
	return "", fmt.Errorf("not implemented")
}

// ValidateToken validates token and returns claims.
func (jc *JWTClient) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: Parse token using jwtlib.ParseWithClaims

	// TODO: Verify signing method matches jc.algo

	// TODO: Return claims if valid or error if invalid/expired
	return nil, fmt.Errorf("not implemented")
}

// ExtractClaims extracts claims from valid token without full validation.
func (jc *JWTClient) ExtractClaims(tokenString string) (*Claims, error) {
	// TODO: Parse token (similar to ValidateToken)

	// TODO: Return claims
	return nil, fmt.Errorf("not implemented")
}
