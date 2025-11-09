package constants

import userEntity "github.com/neurochar/backend/internal/domain/user/entity"

const (
	RightKeyAccessToAdminPanel = "accessToAdminPanel"

	RightKeyAccessToGlobalSettings = "accessToGlobalSettings"
)

// Rights - map of rights
var Rights = map[uint64]*userEntity.Right{
	// доступ в панель управления
	1: {
		ID:                1,
		Key:               RightKeyAccessToAdminPanel,
		Type:              userEntity.RightTypeBool,
		DefaultValue:      0,
		DefaultSuperValue: 1,
	},
	// доступ в глобальные настройки + управление другими пользователями
	2: {
		ID:                2,
		Key:               RightKeyAccessToGlobalSettings,
		Type:              userEntity.RightTypeBool,
		DefaultValue:      0,
		DefaultSuperValue: 1,
	},
}
