package dotastats

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/renstrom/shortuuid"
)

func CreateSessionForUser(email string, mongodb Mongodb) (Session, error) {
	if !govalidator.IsEmail(email) {
		return Session{}, fmt.Errorf("email is not a valid one")
	}
	session := Session{Email: email, SessionKey: shortuuid.New()}
	err := mongodb.CreateOrUpdateSession(session)
	if err != nil {
		return Session{}, fmt.Errorf("session not found, %s", err)
	}

	return session, nil
}

func GetSessionBySessionKey(ssk string, mongodb Mongodb) (Session, error) {
	session, err := mongodb.GetSessionBySessionKey(ssk)
	if err != nil {
		return Session{}, fmt.Errorf("session not found, %s", err)
	}

	return session, nil
}
