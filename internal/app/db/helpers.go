package db

import (
	"aptizer.com/internal/app/models"
	"golang.org/x/crypto/bcrypt"
)

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

func checker(mass []int64, val int64) bool {
	for _, i := range mass {
		if i == val {
			return true
		}
	}
	return false
}

func Reverse(input []*models.Tag) []*models.Tag {
	var output []*models.Tag
	for i := len(input) - 1; i >= 0; i-- {
		output = append(output, input[i])
	}
	return output
}

func remove(slice []*models.News, s int) []*models.News {
	return append(slice[:s], slice[s+1:]...)
}

func tagRemove(slice []*models.Tag, s int64) []*models.Tag {
	return append(slice[:s], slice[s+1:]...)
}

func userRemove(slice []*models.User, s int64) []*models.User {
	return append(slice[:s], slice[s+1:]...)
}

func sliceRemove(slice []int64, s int64) []int64 {
	return append(slice[:s], slice[s+1:]...)
}
