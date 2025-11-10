package parser

import "github.com/gabrielfmcoelho/platform-core/domain"

func ToPublicUserRole(role domain.UserRole) domain.PublicUserRole {
	return domain.PublicUserRole{
		ID:       role.ID,
		RoleName: role.RoleName,
	}
}

func ToPublicUserRoles(roles []domain.UserRole) []domain.PublicUserRole {
	publicRoles := make([]domain.PublicUserRole, len(roles))
	for i, role := range roles {
		publicRoles[i] = ToPublicUserRole(role)
	}
	return publicRoles
}
