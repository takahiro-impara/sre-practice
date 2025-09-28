package service

import (
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password domain.Password) (string, error)
	Compare(hashedPassword domain.Password, plainPassword string) bool
}

type passwordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) PasswordHasher {
	return &passwordHasher{
		cost: cost,
	}
}

func (h *passwordHasher) Hash(password domain.Password) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(string(password)), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h *passwordHasher) Compare(hashedPassword domain.Password, plainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return false
	}
	return true
}
