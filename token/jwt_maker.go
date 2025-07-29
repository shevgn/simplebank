package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeyLength = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, ErrShortSecretKey
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

// CreateToken creates a token with a given duration
func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(j.secretKey))
}

// VerifyToken verifies a token and returns a payload
func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidTokenMethod
		}

		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidTokenPayload
	}

	return payload, nil
}
