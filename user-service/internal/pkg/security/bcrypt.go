package security

import "golang.org/x/crypto/bcrypt"

func HashString(s string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	return hashed, err
}

func CheckStringHash(s string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(s))
}
