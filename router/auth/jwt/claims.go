package jwt

import jwt "github.com/dgrijalva/jwt-go"

// Claims override the default JWT standart claims to add some informations about the user
type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// Valid check if the token is valid (not expired, session key normal, ...)
func (c *Claims) Valid() (err error) {
	if c == nil {
		return ErrClaimsIsNil
	}

	// check time related validity
	if err = c.StandardClaims.Valid(); err != nil {
		return err
	}

	return nil
}
