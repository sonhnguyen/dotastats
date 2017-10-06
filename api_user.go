package dotastats

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const workFactor = 10

func GetUser(email string, pass string, mongodb Mongodb) (User, error) {
	user, err := mongodb.GetUser(email, pass)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateUser(name string, email string, pass string, mongodb Mongodb) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), workFactor)
	if err != nil {
		return User{}, fmt.Errorf("error generating bcrypt hash: %s", err)
	}
	user := User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = mongodb.CreateUser(&user)
	if err != nil {
		return User{}, err
	}

	user.Password = ""

	return user, nil
}
