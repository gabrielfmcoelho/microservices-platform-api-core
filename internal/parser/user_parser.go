package parser

import (
	"github.com/gabrielfmcoelho/platform-core/domain"
)

// Parse User to PublicUser
func ToPublicUser(u domain.User) domain.PublicUser {
	lastLogin := ""
	if len(u.Logs) > 0 {
		// Get the most recent login log
		lastLogin = u.Logs[0].CreatedAt.Format("2006-01-02 15:04:05")
	}

	return domain.PublicUser{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		OrganizationID:   u.Organization.ID,
		OrganizationName: u.Organization.Name,
		RoleID:           u.Role.ID,
		RoleName:         u.Role.RoleName,
		CreatedAt:        u.CreatedAt.Format("2006-01-02 15:04:05"),
		LastLogin:        lastLogin,
		IsArchived:       u.DeletedAt.Valid,
	}
}

// Parse CreateUser to User
func ToUser(cu *domain.CreateUser) *domain.User {
	return &domain.User{
		Name:           cu.Name,
		Email:          cu.Email,
		Password:       cu.Password,
		OrganizationID: cu.OrganizationID,
		RoleID:         cu.RoleID,
	}
}
