package mapper

import (
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
)

func TenantAccountToPb(item *tenantUC.AccountDTO, fileUC fileUC.Usecase, isFullFiles bool) *typesv1.AccountTenant {
	if item == nil {
		return nil
	}

	res := &typesv1.AccountTenant{
		Id:              item.Account.ID.String(),
		Version:         item.Account.Version(),
		RoleId:          item.Account.RoleID,
		Email:           item.Account.Email,
		IsConfirmed:     item.Account.IsConfirmed,
		IsEmailVerified: item.Account.IsEmailVerified,
		IsBlocked:       item.Account.IsBlocked,
		ProfileName:     item.Account.ProfileName,
		ProfileSurname:  item.Account.ProfileSurname,
	}

	if item.Tenant != nil {
		res.TenantId = item.Tenant.ID.String()
	}

	if item.ProfilePhotoOriginalFile != nil && item.ProfilePhoto100x100File != nil {
		res.ProfilePhotos = &typesv1.AccountTenant_PhotoFiles{
			OriginalFile: FileToPb(item.ProfilePhotoOriginalFile, fileUC, isFullFiles),
			S100X100File: FileToPb(item.ProfilePhoto100x100File, fileUC, isFullFiles),
		}
	}

	return res
}
