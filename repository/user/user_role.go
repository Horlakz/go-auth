package user_repository

import (
	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/model"
)

type userRoleRepo struct {
	db database.DatabaseInterface
}

type UserRoleRepositoryInterface interface {
	Create(userRole model.UserRole) (model.UserRole, error)
}

func NewUserRoleRepository(db database.DatabaseInterface) UserRoleRepositoryInterface {
	return &userRoleRepo{db: db}
}

func (u *userRoleRepo) Create(userRole model.UserRole) (model.UserRole, error) {
	userRole.Prepare()

	if err := u.db.Connection().Create(&userRole).Error; err != nil {
		return model.UserRole{}, err
	}

	return userRole, nil
}
