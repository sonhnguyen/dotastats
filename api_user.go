package dotastats

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const workFactor = 10

func GetUserAndAuthenticate(email string, password string, mongodb Mongodb) (User, error) {
	user, err := mongodb.GetUserByEmail(email)
	if err != nil {
		return User{}, fmt.Errorf("user not found, %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, fmt.Errorf("incorrect email or password, %s", err)
	}

	user.Password = ""

	return user, nil
}

func CreateUser(user User, mongodb Mongodb) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), workFactor)
	if err != nil {
		return User{}, fmt.Errorf("error generating bcrypt hash: %s", err)
	}
	user.Password = string(hashedPassword)
	err = mongodb.CreateUser(&user)
	if err != nil {
		return User{}, err
	}

	user.Password = ""

	return user, nil
}
