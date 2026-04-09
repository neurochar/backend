package controller

import (
	"context"
	"net/http"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
	usersV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
)

func (ctrl *Controller) UploadProfilePhotoFile(w http.ResponseWriter, r *http.Request) {
	const op = "UploadProfilePhotoFile"

	ctrl.UploadFileProxy(func(ctx context.Context, file []byte, filename string) ([]byte, string, error) {
		resp, err := ctrl.controls.UsersTenant.UploadProfilePhotoFile(
			ctx,
			&usersV1Pb.UploadProfilePhotoFileRequest{
				File:     file,
				Filename: filename,
			},
		)
		if err != nil {
			return nil, "", appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
		}

		header := server.GatewayMarshaler.ContentType(resp)

		out, err := server.GatewayMarshaler.Marshal(resp)
		if err != nil {
			return nil, "", appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		return out, header, nil
	})(w, r)
}
