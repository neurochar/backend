package roles

import "github.com/neurochar/backend/internal/domain/tenant/entity"

var RoleCreator = entity.Role{
	ID:     1,
	Rank:   1,
	TextID: "creator",
}

var RoleUser = entity.Role{
	ID:     2,
	Rank:   2,
	TextID: "user",
}

var RoleAdmin = entity.Role{
	ID:     3,
	Rank:   3,
	TextID: "admin",
}

var RolesMap = map[uint64]*entity.Role{
	RoleCreator.ID: &RoleCreator,
	RoleUser.ID:    &RoleUser,
	RoleAdmin.ID:   &RoleAdmin,
}
