package user_repository

import (
	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/model"
)

type roleRepository struct {
	db database.DatabaseInterface
}

type RoleRepositoryInterface interface {
	Create(role model.Role) (model.Role, error)
	FindByName(name string) (model.Role, error)
}

func NewRoleRepository(db database.DatabaseInterface) RoleRepositoryInterface {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role model.Role) (model.Role, error) {
	role.Prepare()

	err := r.db.Connection().Create(&role).Error

	if err != nil {
		return model.Role{}, err
	}

	return role, err
}

func (r *roleRepository) FindByName(name string) (model.Role, error) {

	var role model.Role

	err := r.db.Connection().Where("name = ?", name).First(&role).Error

	return role, err
}
