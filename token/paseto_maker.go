package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// Paseto version 2 use Chacha Poly algorithm to encrypt the payload

// PasetoMaker is a paseto token maker
type PasetoMaker struct {
	paseto *paseto.V2
	// use symmetric becs its a local project
	symmetricKey []byte // []byte = field to store the symmetric key
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		// Paseto version 2 use Chacha Poly algorithm to encrypt the payload
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	//else we just create a new PasetoMaker object
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey), // input symmetric Key converted to byte slice
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// else
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil) // nil = the footer
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil) // nil = the footer
	if err != nil {
		return nil, ErrInvalidToken // nil = the payload
	}
	// else check the token is valid or not
	err = payload.Valid()
	if err != nil {
		return nil, err // nil = the payload
	}

	return payload, nil
}
