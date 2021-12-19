package usecases

import (
	"errors"

	"github.com/platelk/contactgraph/domain/models/users"
)

func validateUser(usr *users.User) error {
	if err := validateNickName(usr.NickName); err != nil {
		return err
	}
	return nil
}

func validateNickName(nickName string) error {
	if len(nickName) < 4 || len(nickName) > 20 {
		return errors.New("nickname should have a len greater than 4 and less than 20")
	}
	return nil
}
