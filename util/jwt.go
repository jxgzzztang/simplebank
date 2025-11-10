package util

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	Username string      `json:"username"`
	ID       pgtype.UUID `json:"id"`
	jwt.RegisteredClaims
}

func CreatePayload(username string, duration time.Duration) (TokenPayload, error) {
	// 生成新的UUID
	uuidObj := uuid.New()

	// 转换为pgtype.UUID
	var pgUUID pgtype.UUID
	// 使用正确的方式将uuid.UUID转换为pgtype.UUID
	err := pgUUID.Scan(uuidObj.String())
	if err != nil {
		return TokenPayload{}, err
	}

	payload := TokenPayload{
		Username: username,
		ID:       pgUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    Config.Jwt.Issuer,
		},
	}
	return payload, nil
}

func CreateToken(username string, duration time.Duration) (string, TokenPayload, error) {
	payload, err := CreatePayload(username, duration)
	if err != nil {
		return "", TokenPayload{}, err
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := t.SignedString([]byte(Config.Jwt.SecretKey))
	if err != nil {
		return "", TokenPayload{}, err
	}
	return token, payload, nil
}

func ParseToken(tokenString string) (*TokenPayload, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(Config.Jwt.SecretKey), nil
	})
	if err != nil {
		return nil, false
	}

	if _, ok := token.Claims.(*TokenPayload); !ok || !token.Valid {
		return nil, false
	}

	return token.Claims.(*TokenPayload), true

}
