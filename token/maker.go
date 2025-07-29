// Package token provides a way to handle tokens
package token

import "time"

type Maker interface {
	// CreteToken creates a token with a given duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken verifies a token and returns a payload
	VerifyToken(token string) (*Payload, error)
}
