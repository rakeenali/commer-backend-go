package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// NewHash will inititzalize hash service
func NewHash(salt string) Service {

	return &hashService{
		salt: salt,
	}
}

// Service hash to initialize the service hash
type Service interface {
	GeneratePasswordHash(password string) string
	MatchPassword(hash string, password string) bool
}

type hashService struct {
	salt string
}

func (hs *hashService) GeneratePasswordHash(password string) string {
	pwByte := []byte(password)
	pwByte = append(pwByte, []byte(hs.salt)...)
	hashPw, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	return string(hashPw)
}

func (hs *hashService) MatchPassword(hash string, password string) bool {
	pwByte := []byte(password)
	pwByte = append(pwByte, []byte(hs.salt)...)
	hByte := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hByte, pwByte)
	if err != nil {
		return false
	}
	return true
}
