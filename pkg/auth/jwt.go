package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt" // Import fmt for error messages
	"time"

	"boreholedata-ms/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID utils.BinaryUUID `json:"uid"`
	Role   string           `json:"role"`
	Type   string           `json:"type,omitempty"` // Added token type (access, refresh, reset, email_verify)
	jwt.RegisteredClaims
}

// GenerateAccessToken creates an access token with a specific type claim.
func GenerateAccessToken(
	userID utils.BinaryUUID,
	role string,
	secret string,
	expiry time.Duration,
) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   "access", // Explicitly mark as access token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken creates a refresh token with a specific type claim and a longer expiry.
// Role is generally not included in refresh tokens.
func GenerateRefreshToken(
	userID utils.BinaryUUID,
	secret string, // Use a separate secret for refresh tokens, or derived
	expiry time.Duration,
) (string, error) {
	claims := Claims{
		UserID: userID,
		// Role is typically not needed in refresh tokens.
		Type: "refresh", // Explicitly mark as refresh token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GeneratePasswordResetToken creates a password reset token.
func GeneratePasswordResetToken(
	userID utils.BinaryUUID,
	secret string,
	expiry time.Duration,
) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   "reset", // Explicitly mark as reset token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateEmailVerificationToken creates an email verification token.
func GenerateEmailVerificationToken(
	userID utils.BinaryUUID,
	secret string,
	expiry time.Duration,
) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   "email_verify", // Explicitly mark as email verification token
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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]) // More specific error
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims // Use a specific JWT error if parsing failed but token was generally valid
}

// HashToken creates a secure SHA256 hash for storage (e.g., for refresh tokens, password reset tokens)
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// FormatTokenResponse prepares tokens for HTTP response (Removed from here, will be handled in service/handler DTO)
// func FormatTokenResponse(accessToken, refreshToken string) string {
// 	return "Bearer " + accessToken + "|" + refreshToken
// }
