package seed

import (
	"fmt"

	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/model"
)

type SeederInterface interface {
	Seed()
}

type seeder struct {
	dbConn database.DatabaseInterface
}

func NewSeeder(dbConn database.DatabaseInterface) SeederInterface {
	return &seeder{dbConn: dbConn}
}

func (s *seeder) Seed() {
	s.SeedRoles()
}

func (s *seeder) SeedRoles() {
	roles := []model.Role{
		{Name: model.ROLE_CORPORATE},
		{Name: model.ROLE_INDIVIDUAL},
		{Name: model.ROLE_JOINT},
	}

	for _, role := range roles {
		roleExists := s.dbConn.Connection().Where("name = ?", role.Name).First(&model.Role{}).RowsAffected > 0
		if roleExists {
			fmt.Printf("%s role already exists in the database. Skipping seeding...\n", role.Name)
		} else {
			role.Prepare()
			if err := s.dbConn.Connection().Create(&role).Error; err != nil {
				fmt.Println("Failed to create role:", err)
			}
			fmt.Println(role.Name, "role created successfully.")
		}
	}

}
