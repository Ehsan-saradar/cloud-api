package auth

import (
	"context"
	"time"

	"api.cloud.io/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Tokens life time
const (
	RefreshTokenLifeTime = time.Hour * 24 * 365
	AccessTokenLifeTime  = time.Minute * 10
)

var (
	jwtSigningKey = []byte{108, 40, 37, 154, 163, 13, 79, 29,
		250, 220, 81, 145, 212, 91, 73, 41,
		195, 135, 117, 254, 254, 42, 184, 251,
		62, 90, 236, 212, 49, 105, 41, 27}

	jwtSigner, _ = jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       jwtSigningKey},
		(&jose.SignerOptions{}).WithType("JWT"))
)

// Token is the model of access/refresh token in JWT fromat.
type Token struct {
	ID             uuid.UUID `json:"jti"`
	Subject        uuid.UUID `json:"sub"`
	SubjectType    string    `json:"sub_typ"`
	IssuedAt       int64     `json:"iat"`
	ExpirationTime int64     `json:"exp"`
	Scopes         []string  `json:"scopes"`
}

// ContextToken get the Token of request.
func ContextToken(ctx context.Context) *Token {
	if ctx.Value(tokenContextKey) != nil {
		return ctx.Value(tokenContextKey).(*Token)
	}

	return nil
}

// ParseSignedToken unmarshal and verify signature of token.
func ParseSignedToken(str string) (*Token, error) {
	rawJWT, err := jwt.ParseSigned(str)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	var token Token
	err = rawJWT.Claims(jwtSigningKey, &token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	return &token, nil
}

// Validate validate the token.ExpirationTime against the now time.
func (token *Token) Validate() error {
	if token.ExpirationTime < time.Now().Unix() {
		return errors.ErrTokenExpired
	}

	return nil
}

// SignedString returns JWS of token.
func (token *Token) SignedString() (string, error) {
	str, err := jwt.Signed(jwtSigner).Claims(token).CompactSerialize()
	if err != nil {
		return "", errors.ErrInternal(err)
	}

	return str, nil
}
