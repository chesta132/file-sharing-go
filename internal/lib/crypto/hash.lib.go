package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	pw := []byte(password)
	cost := bcrypt.DefaultCost

	hpw, _ := bcrypt.GenerateFromPassword(pw, cost)
	return string(hpw)
}

func ComparePassword(hashed, password string) bool {
	pw := []byte(password)
	hpw := []byte(hashed)

	err := bcrypt.CompareHashAndPassword(hpw, pw)
	return err == nil
}
