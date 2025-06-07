package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"boreholedata-ms/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID utils.BinaryUUID `json:"uid"`
	Role   string           `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates access and refresh tokens
func GenerateToken(
	userID utils.BinaryUUID,
	role string,
	secret string,
	expiry time.Duration,
) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateAccessToken creates only an access token
func GenerateAccessToken(
	userID utils.BinaryUUID,
	role string,
	secret string,
	expiry time.Duration,
) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken verifies token and returns claims
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// HashToken creates a secure hash for storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// FormatTokenResponse prepares tokens for HTTP response
func FormatTokenResponse(accessToken, refreshToken string) string {
	return "Bearer " + accessToken + "|" + refreshToken
}
