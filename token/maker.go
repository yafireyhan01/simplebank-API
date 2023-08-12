package token

import "time"

// Maker is an interface for managing token
type Maker interface {
	// CreateToken creates a new token for a specific username and duratoin
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not, if valid it will return Payload data
	// stored inside the body of the token
	VerifyToken(token string) (*Payload, error)
}
