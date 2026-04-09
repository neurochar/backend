package crm_tenant

import (
	"context"
	"fmt"

	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
)

func (ctrl *controller) UploadCandidateResumeFile(
	ctx context.Context,
	req *desc.UploadCandidateResumeFileRequest,
) (*desc.UploadCandidateResumeFileResponse, error) {
	ctx = auth.WithCheckTenantAccess(ctx)

	userAuthData := auth.GetAuthData(ctx)

	fmt.Println(2, userAuthData)

	fmt.Println(req.File)

	return &desc.UploadCandidateResumeFileResponse{
		FileId: "GGWP",
	}, nil
}
