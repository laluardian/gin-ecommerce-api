package libs

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw *string) error {
	if len(*pw) == 0 {
		return errors.New("password should not be empty")
	}

	bytePw := []byte(*pw)
	bytes, err := bcrypt.GenerateFromPassword(bytePw, 10)
	if err != nil {
		return err
	}

	*pw = string(bytes)

	return nil
}

func ComparePassword(hash, pw string) bool {
	bytePw := []byte(pw)
	byteHash := []byte(hash)
	pwCompareResult := bcrypt.CompareHashAndPassword(byteHash, bytePw) == nil

	return pwCompareResult
}
