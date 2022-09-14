package db

import "golang.org/x/crypto/bcrypt"

var bcryptRounds = 10

// ComparePassword - Necessary to verify the authenticity of the password during login to the account.
func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// HashPassword - Necessary for hashing the user's password.
func (m *UsersStorage) HashPassword(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcryptRounds)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
