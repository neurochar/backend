package users_tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) UpdateMyProfile(
	ctx context.Context,
	req *desc.UpdateMyProfileRequest,
) (*desc.UpdateMyProfileResponse, error) {
	const op = "UpdateMyProfile"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithoutCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	var profilePhotos *tenantUC.AccountDataInputProfilePhotos

	if req.Payload.GetProfilePhotos() != nil {
		profilePhotos = &tenantUC.AccountDataInputProfilePhotos{}

		parseID, err := uuid.Parse(req.Payload.GetProfilePhotos().OriginalFileId)
		if err != nil {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		profilePhotos.PhotoOriginalFileID = lo.ToPtr(parseID)

		parseID, err = uuid.Parse(req.Payload.GetProfilePhotos().S100X100FileId)
		if err != nil {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		profilePhotos.Photo100x100FileID = lo.ToPtr(parseID)

	} else if req.Payload.GetProfilePhotosClear() {
		profilePhotos = &tenantUC.AccountDataInputProfilePhotos{}
	}

	usecaseInput := tenantUC.PatchAccountDataInput{
		Version: req.Version,

		ProfileName:    &req.Payload.ProfileName,
		ProfileSurname: &req.Payload.ProfileSurname,
		ProfilePhotos:  profilePhotos,
	}

	err := ctrl.tenantFacade.Account.PatchAccountByDTO(
		ctx,
		authData.TenantUserClaims().AccountID,
		usecaseInput,
		req.SkipVersionCheck,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.UpdateMyProfileResponse{}, nil
}
