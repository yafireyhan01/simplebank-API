package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// length key for better security (should be too short) 32 character
const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	// use symmetric key
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be atleast %d character", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil // nil = no error
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// get signing algorithm via token.method field
		_, ok := token.Method.(*jwt.SigningMethodHMAC) // conver to HMAC bcs using HS256, which is the instance of HMAC struct
		if !ok {
			return nil, ErrInvalidToken
		}
		// if conversion was successful (the algorithm matches)
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// verr = validaiton error
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	// incase everything is good
	// and the token is successfully parse and verified
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
