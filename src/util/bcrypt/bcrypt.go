package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func ValidatePwd(pwd string, pwdHash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd)); err != nil {
		return err
	}
	return nil
}

func PwdHash(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}
