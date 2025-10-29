package util

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production")

// Claims represents JWT claims
type Claims struct {
	UserID int32 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID int32, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates and parses JWT token
func ValidateJWT(tokenString string) (int32, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}

// Generate2FASecret generates a new TOTP secret
func Generate2FASecret() string {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(secret)
}

// Generate2FAQRCode generates QR code URL for 2FA setup
func Generate2FAQRCode(secret, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", issuer, secret, issuer)
}

// Verify2FACode verifies TOTP code
func Verify2FACode(secret, code string) bool {
	valid := totp.Validate(code, secret)
	return valid
}

// GenerateBackupCodes generates backup codes for 2FA
func GenerateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 4)
		rand.Read(b)
		codes[i] = fmt.Sprintf("%X-%X", b[:2], b[2:])
	}
	return codes
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// HashPassword is a helper for bcrypt hashing (imported from golang.org/x/crypto/bcrypt)
// Use: bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	// Simple email validation
	// In production, use a proper email validation library
	return len(email) > 3 && len(email) < 255
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	// Add more validation rules as needed
	return nil
}
