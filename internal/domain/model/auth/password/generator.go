package password_auth

import (
	"github.com/sethvargo/go-password/password"
)

func GeneratePassword(length int) (string, error) {
	pass, err := password.Generate(length, 10, 10, true, false)
	if err != nil {
		return "", err
	}
	return pass, nil
}
