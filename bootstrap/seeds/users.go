// bootstrap/seeds/user_seed.go
package seeds

import (
	"log"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
	"github.com/gabrielfmcoelho/platform-core/internal/password"
	"gorm.io/gorm"
)

// SeedUsers popula a tabela de usuários com dados iniciais
func SeedUsers(db *gorm.DB) error {
	// Exemplo: Verificar se já existe algum usuário
	var count int64
	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Se não tiver usuário nenhum, cria alguns
	if count == 0 {
		users := []domain.CreateUser{
			{
				Name:           "Solude Admin",
				Email:          "contato@solude.tech",
				Password:       "admin123",
				OrganizationID: 1,
				RoleID:         1, // Ex: Role Admin
			},
			{
				Name:           "Gabriel Coelho",
				Email:          "gabrielcoelho@inovadata.tech",
				Password:       "123",
				OrganizationID: 1,
				RoleID:         1, // Ex: Role User
			},
			{
				Name:           "HSM Contact",
				Email:          "contato@hsm.com",
				Password:       "123",
				OrganizationID: 2,
				RoleID:         2, // Ex: Role Manager
			},
			{
				Name:           "Guest User",
				Email:          "guest@solude.tech",
				Password:       "123",
				OrganizationID: 4,
				RoleID:         4, // Ex: Role Guest
			},
		}

		for _, u := range users {
			hashedPassword, err := password.HashPassword(u.Password)
			if err != nil {
				return err
			}
			u.Password = hashedPassword
			user := parser.ToUser(&u)
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
		log.Printf("Users seed executado: criados %d usuários\n", len(users))
	}
	return nil
}
