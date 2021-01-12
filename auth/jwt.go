package auth

import (
	"commerce/helpers"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User is a payload for jwt
type User struct {
	Username string
	ID       uint
}

type claims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}

// InitAuth will initialize and setup auth
func InitAuth(secret string) Auth {
	return &authService{
		secret: secret,
	}
}

// Auth interface that any auth service must implements
type Auth interface {
	SignToken(user *User) string
	VerifyToken(token string) (*User, error)
}

type authService struct {
	secret string
}

func (as *authService) SignToken(user *User) string {
	c := claims{
		Username: user.Username,
		ID:       user.ID,
	}
	c.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()
	c.Issuer = "localhost"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	jwt, err := token.SignedString([]byte(as.secret))
	if err != nil {
		panic(err)
	}
	return jwt
}

func (as *authService) VerifyToken(token string) (*User, error) {
	var c claims
	t, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.secret), nil
	})

	if err != nil || !t.Valid {
		return nil, helpers.ErrInvalidToken
	}

	return &User{
		Username: c.Username,
		ID:       c.ID,
	}, nil
}
