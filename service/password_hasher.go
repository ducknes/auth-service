package service

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Hasher struct {
	PasswordSalt string
}

func NewHasher(passwordSalt string) PasswordHasher {
	return &Hasher{
		PasswordSalt: passwordSalt,
	}
}

func (h *Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.PasswordSalt))), nil
}
