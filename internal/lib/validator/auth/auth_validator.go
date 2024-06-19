package auth

import (
	"sso/internal/lib/validator"
)

func ValidateLoginReq(email, password string, appId int32) error {
	if err := validator.ValidateEmail(email); err != nil {
		return err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return err
	}
	if err := validator.ValidateAppId(appId); err != nil {
		return err
	}
	return nil
}

func ValidateRegisterReq(email, password string) error {
	if err := validator.ValidateEmail(email); err != nil {
		return err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return err
	}
	return nil
}
