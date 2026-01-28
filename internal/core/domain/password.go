package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomPassword() (string, error) {
	const passwordLength = 6
	password := make([]byte, passwordLength)

	for i := 0; i < passwordLength; i++ {
		// Gera um número aleatório entre 0 e 9
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("erro ao gerar número aleatório: %w", err)
		}
		password[i] = byte('0' + num.Int64())
	}

	return string(password), nil
}

// HashPassword makes a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}
	return string(bytes), nil
}

// CheckPasswordHash verifies if the password matches the hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
