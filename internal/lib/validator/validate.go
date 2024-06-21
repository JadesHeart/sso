package validator

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

const (
	EmailInvalid    = "invalid email"
	PasswordInvalid = "invalid password"
	AppIdInvalid    = "invalid app_id"
)

const (
	emptyValue = 0
)

func ValidateEmail(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return status.Error(codes.InvalidArgument, EmailInvalid)
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(emailRegex, str)

	if !match {
		return status.Error(codes.InvalidArgument, EmailInvalid)
	}
	return nil
}

func ValidatePassword(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return status.Error(codes.InvalidArgument, PasswordInvalid)
	}

	if len(str) > 20 {
		return status.Error(codes.InvalidArgument, PasswordInvalid)
	}

	passwordRegex := `^[a-zA-Z0-9!?\.&]+$`

	match, _ := regexp.MatchString(passwordRegex, str)

	if !match {
		return status.Error(codes.InvalidArgument, PasswordInvalid)
	}

	return nil
}

func ValidateAppId(value interface{}) error {
	appId, ok := value.(int32)
	if !ok {
		return status.Error(codes.InvalidArgument, AppIdInvalid)
	}
	if appId == emptyValue {
		return status.Error(codes.InvalidArgument, AppIdInvalid)
	}

	return nil
}
