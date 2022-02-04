package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

func UserGeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hashedPassword), err
}

func UserVerifyPassword(userPwd, userStoredPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userStoredPwd), []byte(userPwd))
	return err == nil
}
