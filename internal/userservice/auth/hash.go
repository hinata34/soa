//go:generate mockgen -source=./hash.go -destination=./mocks/hash.go -package=mock
package auth

import "golang.org/x/crypto/bcrypt"

type Hash interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type HashImplementation struct {
}

func (h *HashImplementation) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
