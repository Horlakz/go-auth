package service

import (
	"errors"
	"time"

	"github.com/horlakz/go-auth/dto"
	"github.com/horlakz/go-auth/helper"
	"github.com/horlakz/go-auth/model"
	user_repository "github.com/horlakz/go-auth/repository/user"
	"gorm.io/gorm"
)

type authService struct {
	userRepo     user_repository.UserRepositoryInterface
	roleRepo     user_repository.RoleRepositoryInterface
	userRoleRepo user_repository.UserRoleRepositoryInterface
	codeRepo     user_repository.CodeRepositoryInterface
	hashing      helper.HashingInterface
	jwt          helper.JwtInterface
	smtp         SmtpInterface
}

type AuthServiceInterface interface {
	Login(loginDto dto.LoginDto) (dto.LoginResponseDTO, error)
	RefreshAccessToken(token string) (dto.LoginResponseDTO, error)
	Register(registerDto dto.RegisterDto) error
	ForgotPassword(email string) error
	ResetPassword(resetPasswordDto dto.ResetPasswordDto) error
}

func NewAuthService(
	userRepo user_repository.UserRepositoryInterface,
	roleRepo user_repository.RoleRepositoryInterface,
	userRoleRepo user_repository.UserRoleRepositoryInterface,
	codeRepo user_repository.CodeRepositoryInterface,
	smtp SmtpInterface,
) AuthServiceInterface {
	return &authService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		codeRepo:     codeRepo,
		hashing:      helper.NewHashing(),
		jwt:          helper.NewJwt(),
		smtp:         smtp,
	}
}

// Login implements AuthServiceInterface.
func (a *authService) Login(loginDto dto.LoginDto) (dto.LoginResponseDTO, error) {
	user, err := a.userRepo.FindUserByEmail(loginDto.Email)

	if err == gorm.ErrRecordNotFound {
		return dto.LoginResponseDTO{}, errors.New("user does not exist")
	}

	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	match, err := a.hashing.ComparePassword(loginDto.Password, user.Password)

	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	if !match {
		return dto.LoginResponseDTO{}, errors.New("invalid password")
	}

	return a.CreateTokens(user.ID.String()), nil

}

// RefreshAccessToken implements AuthServiceInterface.
func (a *authService) RefreshAccessToken(token string) (dto.LoginResponseDTO, error) {
	userId, err := a.jwt.ExtractUserID(token, "refresh")

	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	return a.CreateTokens(userId.String()), nil

}

// Register implements AuthServiceInterface.
func (a *authService) Register(registerDto dto.RegisterDto) error {
	var user model.User
	var userRole model.UserRole

	hash, err := a.hashing.HashPassword(registerDto.Password)

	if err != nil {
		return err
	}

	user.Fullname = registerDto.FullName
	user.Email = registerDto.Email
	user.Password = hash

	if _, err = a.userRepo.FindUserByEmail(registerDto.Email); err != gorm.ErrRecordNotFound {
		return errors.New("user already exists")
	}

	role, err := a.roleRepo.FindByName(registerDto.Role)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("role does not exist")
		}

		return err
	}

	newUser, err := a.userRepo.Create(user)

	if err != nil {
		// a.userRepo.DeleteUser(newUser.ID)
		return err
	}

	userRole.UserID = newUser.ID
	userRole.RoleID = role.ID

	_, err = a.userRoleRepo.Create(userRole)

	if err != nil {
		return err
	}

	return nil

}

// ForgotPassword implements AuthServiceInterface.
func (a *authService) ForgotPassword(email string) error {
	var sendEmailDto dto.SendEmailDto
	var verifCode model.VerificationCode

	user, err := a.userRepo.FindUserByEmail(email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user does not exist")
		}

		return err
	}

	code := helper.GenerateRandomDigits(6)

	verifCode.UserID = user.ID
	verifCode.Code = code
	verifCode.Purpose = "forgot_password"

	if _, err = a.codeRepo.Create(verifCode); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"FullName": user.Fullname,
		"Code":     code,
	}

	sendEmailDto.Subject = "Reset Password Request"
	sendEmailDto.To = user.Email
	sendEmailDto.Template = "reset-password"
	sendEmailDto.Variables = vars

	go func(sendEmailDto dto.SendEmailDto) {
		a.smtp.SendWithTemplate(sendEmailDto)
	}(sendEmailDto)

	return nil
}

// ResetPassword implements AuthServiceInterface.
func (a *authService) ResetPassword(resetPasswordDto dto.ResetPasswordDto) error {
	verifCode, err := a.codeRepo.FindByCodeAndEmail(resetPasswordDto.Code, resetPasswordDto.Email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("email or code does not exists")
		}

		return err
	}

	if time.Now().After(verifCode.CreatedAt.Add(10 * time.Minute)) {
		return errors.New("code has expired")
	}

	hash, err := a.hashing.HashPassword(resetPasswordDto.Password)

	if err != nil {
		return err
	}

	user, err := a.userRepo.FindUserByEmail(resetPasswordDto.Email)

	if err != nil {
		return err
	}

	user.Password = hash

	if _, err = a.userRepo.Update(user); err != nil {
		return err
	}

	if err = a.codeRepo.Delete(verifCode); err != nil {
		return err
	}

	return nil
}

func (a *authService) CreateTokens(userId string) dto.LoginResponseDTO {
	accessToken, err := a.jwt.CreateToken(userId, "access")

	if err != nil {
		return dto.LoginResponseDTO{}
	}

	refreshToken, err := a.jwt.CreateToken(userId, "refresh")

	if err != nil {
		return dto.LoginResponseDTO{}
	}

	tokenDto := dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenDto
}
