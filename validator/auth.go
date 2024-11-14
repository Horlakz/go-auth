package validator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/horlakz/go-auth/payload/request"
)

type AuthValidator struct {
	Validator[request.RegisterRequest]
}

func (validator *AuthValidator) LoginValidate(loginReq request.LoginRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&loginReq,
		validation.Field(&loginReq.Email, validation.Required, is.Email, validation.Length(3, 32)),
		validation.Field(&loginReq.Password, validation.Required, validation.Length(3, 32)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *AuthValidator) RegisterValidate(registerDto request.RegisterRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&registerDto,
		validation.Field(&registerDto.FullName, validation.Required, validation.Length(3, 256)),
		validation.Field(&registerDto.Email, validation.Required, is.Email, validation.Length(3, 32)),
		validation.Field(&registerDto.Password, validation.Required, validation.Length(3, 32)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *AuthValidator) EmailValidate(emailReq request.EmailRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&emailReq,
		validation.Field(&emailReq.Email, validation.Required, is.Email, validation.Length(3, 32), is.Email),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *AuthValidator) ResetPasswordValidate(resetPasswordReq request.ResetPasswordRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&resetPasswordReq,
		validation.Field(&resetPasswordReq.Email, validation.Required, is.Email, validation.Length(3, 32)),
		validation.Field(&resetPasswordReq.Password, validation.Required, validation.Length(3, 32)),
		validation.Field(&resetPasswordReq.Code, validation.Required, validation.Length(6, 6)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
