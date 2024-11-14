package user_repository

import (
	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/model"
)

type codeRepo struct {
	db database.DatabaseInterface
}

type CodeRepositoryInterface interface {
	Create(code model.VerificationCode) (model.VerificationCode, error)
	FindByCodeAndEmail(code, email string) (model.VerificationCode, error)
	Delete(code model.VerificationCode) error
}

func NewVerificationCodeRepository(db database.DatabaseInterface) CodeRepositoryInterface {
	return &codeRepo{db: db}
}

func (u *codeRepo) Create(code model.VerificationCode) (model.VerificationCode, error) {
	code.Prepare()

	if err := u.db.Connection().Create(&code).Error; err != nil {
		return model.VerificationCode{}, err
	}

	return code, nil
}

func (u *codeRepo) FindByCodeAndEmail(code, email string) (model.VerificationCode, error) {
	var verificationCode model.VerificationCode

	if err := u.db.Connection().Joins("JOIN users ON users.id = verification_codes.user_id").Where("verification_codes.code = ? AND users.email = ?", code, email).First(&verificationCode).Error; err != nil {
		return model.VerificationCode{}, err
	}

	return verificationCode, nil
}

func (u *codeRepo) Delete(code model.VerificationCode) error {
	if err := u.db.Connection().Delete(&code).Error; err != nil {
		return err
	}

	return nil
}
