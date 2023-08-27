package token

import "time"

// Maker is an interface for managing token
type Maker interface {
	// CreateToken create a new token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// ValidateToken check if the token valid or not
	ValidateToken(token string) (*Payload, error)
}
