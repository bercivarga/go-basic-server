package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWT_DURATION         = 24 * time.Hour * 7  // 7 days
	JWT_REFRESH_DURATION = 24 * time.Hour * 14 // 14 days
)

type JWTManager struct {
	secretKey       string
	TokenDuration   time.Duration
	RefreshDuration time.Duration
}

type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey, JWT_DURATION, JWT_REFRESH_DURATION}
}

func (j *JWTManager) Generate(userID int64) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) Verify(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func (j *JWTManager) CreateExpiry() (accessTokenExpireAt, refreshTokenExpireAt time.Time) {
	accessTokenExpireAt = time.Now().Add(j.TokenDuration)
	refreshTokenExpireAt = time.Now().Add(j.RefreshDuration)
	return
}
