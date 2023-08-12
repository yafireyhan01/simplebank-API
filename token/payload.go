package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID // to invalidate some specific token in case they are leaked
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`  // when the token is created
	ExpiredAt time.Time `json:"expired_at"` // when the token is expired
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	// if not error
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// valid checks, if the token payload is valid or not
func (payload *Payload) Valid() error {
	// check the expiration time of the token
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	// else if the token not expired
	return nil
}
