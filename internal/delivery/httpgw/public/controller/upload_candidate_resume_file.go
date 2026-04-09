package controller

import (
	"context"
	"net/http"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
	crmV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
)

func (ctrl *Controller) UploadCandidateResumeFile(w http.ResponseWriter, r *http.Request) {
	const op = "UploadCandidateResumeFile"

	ctrl.UploadFileProxy(func(ctx context.Context, file []byte, filename string) ([]byte, string, error) {
		resp, err := ctrl.controls.CRM.UploadCandidateResumeFile(
			ctx,
			&crmV1Pb.UploadCandidateResumeFileRequest{
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
