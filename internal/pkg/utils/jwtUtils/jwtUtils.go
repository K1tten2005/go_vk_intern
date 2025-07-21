package jwtUtils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth"
	"github.com/golang-jwt/jwt"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/require"
)

type CtxKey string

const (
	UserIdKey CtxKey = "id"
	UserLoginKey CtxKey = "login"
)


func GenerateToken(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", auth.ErrGeneratingToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func GetIdFromJWT(JWTStr string, secret string) (string, bool) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(JWTStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET is not set")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", false
	}

	id, ok := claims["id"].(string)
	return id, ok
}

func GetLoginFromJWT(JWTStr string, secret string) (string, bool) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(JWTStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET is not set")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", false
	}

	login, ok := claims["login"].(string)
	return login, ok
}

func GetIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, err := uuid.FromString(ctx.Value(UserIdKey).(string))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func GetLoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(UserLoginKey).(string)
	return login, ok
}

func GenerateJWTForTest(t *testing.T, login string, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":  login,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenStr
}
