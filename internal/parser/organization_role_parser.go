package parser

import (
	"github.com/gabrielfmcoelho/platform-core/domain"
)

func ToPublicOrganizationRole(role domain.OrganizationRole) domain.PublicOrganizationRole {
	return domain.PublicOrganizationRole{
		ID:       role.ID,
		RoleName: role.RoleName,
	}
}

func ToPublicOrganizationRoles(roles []domain.OrganizationRole) []domain.PublicOrganizationRole {
	publicRoles := make([]domain.PublicOrganizationRole, len(roles))
	for i, role := range roles {
		publicRoles[i] = ToPublicOrganizationRole(role)
	}
	return publicRoles
}
