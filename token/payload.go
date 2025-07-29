package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken        = errors.New("token is expired")
	ErrInvalidTokenPayload = errors.New("invalid token payload")
	ErrInvalidTokenMethod  = errors.New("invalid token method")
	ErrShortSecretKey      = fmt.Errorf("secret key must be at least %d characters long", minSecretKeyLength)
)

// Payload is a token payload
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) *Payload {
	tokenID := uuid.New()

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload
}

// func (p *Payload) Valid() error {
// 	if time.Now().After(p.ExpireAt) {
// 		return ErrExpiredToken
// 	}
//
// 	return nil
// }

// GetAudience implements jwt.Claims.
func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}

// GetExpirationTime implements jwt.Claims.
func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

// GetIssuedAt implements jwt.Claims.
func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

// GetIssuer implements jwt.Claims.
func (p *Payload) GetIssuer() (string, error) {
	return "simplebank", nil
}

// GetNotBefore implements jwt.Claims.
func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

// GetSubject implements jwt.Claims.
func (p *Payload) GetSubject() (string, error) {
	return p.Username, nil
}
