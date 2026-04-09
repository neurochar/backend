package users_tenant

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	"github.com/neurochar/backend/pkg/auth"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
)

func (ctrl *Controller) UploadProfilePhotoFile(
	ctx context.Context,
	req *desc.UploadProfilePhotoFileRequest,
) (*desc.UploadProfilePhotoFileResponse, error) {
	const op = "UploadProfilePhotoFile"

	ip := tools.GetRealIP(ctx)

	err := ctrl.limiter.Get(limiter.DefaultName).Register(ctx, &limiter.RegisterKey{
		IP: ip,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	files, err := ctrl.tenantFacade.Account.UploadProfileImageFile(ctx, req.Filename, req.File)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.UploadProfilePhotoFileResponse{
		Data: &typesv1.FilesMap{
			Map: mapper.FilesToMapPb(files, ctrl.fileUC, true),
		},
	}, nil
}
