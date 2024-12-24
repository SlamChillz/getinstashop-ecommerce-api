package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	MininumAllowedSecretKeySize = 32
	TokenDuration               = time.Hour * 12
	ErrTokenIsInvalid           = errors.New("token is invalid")
)

type Payload struct {
	UserID               uuid.UUID `json:"userId"`
	Admin                bool      `json:"admin"`
	jwt.RegisteredClaims `json:"claims"`
}

func (p *Payload) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func NewPayload(userId uuid.UUID, admin bool) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		UserID: userId,
		Admin:  admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
		},
	}
	return payload, err
}

type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) (*JWT, error) {
	if len(secretKey) < MininumAllowedSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: key must be at least %d characters", MininumAllowedSecretKeySize)
	}
	return &JWT{secretKey}, nil
}

func (jwtToken *JWT) CreateToken(userId uuid.UUID, admin bool) (string, error) {
	payload, err := NewPayload(userId, admin)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(jwtToken.secretKey))
	return tokenString, err
}

func (jwtToken *JWT) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method: %v", token.Header["alg"])
		}
		return []byte(jwtToken.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		// Log the error
		return nil, ErrTokenIsInvalid
	}
	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrTokenIsInvalid
	}
	return payload, nil
}
