package user_repository

import (
	"github.com/google/uuid"
	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/model"
)

type userRepo struct {
	db database.DatabaseInterface
}

type UserRepositoryInterface interface {
	Create(user model.User) (model.User, error)
	FindUserByEmail(email string) (model.User, error)
	Update(user model.User) (model.User, error)
}

func NewUserRepository(db database.DatabaseInterface) UserRepositoryInterface {
	return &userRepo{db: db}
}

func (ur *userRepo) Create(user model.User) (model.User, error) {
	user.Prepare()

	err := ur.db.Connection().Create(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, err
}

func (ur *userRepo) FindUserByEmail(email string) (model.User, error) {
	var user model.User

	err := ur.db.Connection().Where("email = ?", email).First(&user).Error

	return user, err
}

func (ur *userRepo) FindById(id uuid.UUID) (model.User, error) {
	var user model.User

	err := ur.db.Connection().Where("id = ?", id).First(&user).Error

	return user, err
}

func (ur *userRepo) Update(user model.User) (model.User, error) {

	err := ur.db.Connection().Save(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, err
}
