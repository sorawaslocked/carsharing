package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return hashed, err
}

func CheckPassword(password string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
