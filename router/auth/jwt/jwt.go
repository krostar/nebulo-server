package jwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/router/auth"
)

var (
	// ErrClaimsIsNil is an error returned
	// when the claims object is nil
	ErrClaimsIsNil = errors.New("claims is nil")
)

// NewToken generate and sign a new JWT token
func NewToken(JWTSecret string, expire int64, user string) (t *auth.Token, err error) {
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expire,
		},
	}

	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := JWTToken.SignedString([]byte(JWTSecret))
	if err != nil {
		log.Errorln("JWT: unable to sign")
		return nil, err
	}
	return &auth.Token{
		Token:  token,
		Expire: expire,
	}, nil
}

// GetUser retrieve the user stored in the JWT claims
func GetUser(key interface{}) (string, error) {
	token, ok := key.(*jwt.Token)
	if !ok {
		return "", errors.New("JWT: bad token type")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("JWT: bad claim type")
	}

	if err := claims.Valid(); err != nil {
		return "", err
	}

	return claims.User, nil
}
