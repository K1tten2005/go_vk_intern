package jwtUtils

import (
	"testing"
	"time"

	"github.com/satori/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

const secret = "test_secret"

func createTestJWT(t *testing.T, claims jwt.MapClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	return tokenStr
}

func TestGetLoginFromJWT(t *testing.T) {
	claims := jwt.MapClaims{
		"login": "test2025",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	tokenStr := createTestJWT(t, claims, secret)

	login, ok := GetLoginFromJWT(tokenStr, secret)
	assert.True(t, ok)
	assert.Equal(t, "test2025", login)
}

func TestGetIdFromJWT(t *testing.T) {
	userId := uuid.NewV4()
	claims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	tokenStr := createTestJWT(t, claims, secret)

	idStr, ok := GetIdFromJWT(tokenStr, secret)
	id := uuid.FromStringOrNil(idStr)
	assert.True(t, ok)
	assert.Equal(t, userId, id)
}

func TestGenerateJWTForTest(t *testing.T) {
	login := "test2025"
	secret := "secret"

	tokenStr := GenerateJWTForTest(t, login, secret)
	assert.NotNil(t, tokenStr)
}
