package service

import (
	"context"

	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

type RoleService struct {
	storage *postgresql.Storage
}

func (s *RoleService) GetRole(ctx context.Context, roleName string) (*postgresql.Role, error) {
	role, err := s.storage.Roles.GetByName(ctx, roleName)
	if err != nil {
		return nil, err
	}

	return role, nil
}
