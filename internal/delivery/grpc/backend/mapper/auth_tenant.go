package mapper

import (
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	auth_tenantv1 "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func AuthTenantAccountToPb(item *tenantUC.AccountDTO, fileUC fileUC.Usecase, isFullFiles bool) *auth_tenantv1.Account {
	if item == nil {
		return nil
	}

	return &auth_tenantv1.Account{
		Id:                       item.Account.ID.String(),
		Version:                  item.Account.Version(),
		RoleId:                   item.Account.RoleID,
		Email:                    item.Account.Email,
		IsConfirmed:              item.Account.IsConfirmed,
		IsEmailVerified:          item.Account.IsEmailVerified,
		ProfileName:              item.Account.ProfileName,
		ProfileSurname:           item.Account.ProfileSurname,
		ProfilePhotoOriginalFile: FileToPb(item.ProfilePhotoOriginalFile, fileUC, isFullFiles),
		ProfilePhoto_100X100File: FileToPb(item.ProfilePhoto100x100File, fileUC, isFullFiles),
	}
}
