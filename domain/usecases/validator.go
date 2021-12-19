package usecases

import (
	"errors"

	"github.com/platelk/contactgraph/domain/models/users"
)

func validateUser(usr *users.User) error {
	if err := validateNickName(usr.NickName); err != nil {
		return err
	}
	if err := validatePhoneNumber(string(usr.PhoneNumber)); err != nil {
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

func validatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) == 0 {
		return errors.New("phone number required")
	}
	if phoneNumber[0] != '+' {
		return errors.New("phone number need to be in format '+1 123 456 789 00'")
	}
	if len(phoneNumber) > 18 {
		return errors.New("phone number is too large")
	}
	return nil
}
