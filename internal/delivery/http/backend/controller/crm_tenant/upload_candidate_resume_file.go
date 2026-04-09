package crm_tenant

import (
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/dto"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
)

func (ctrl *Controller) UploadCandidateResumeFileHandler(c *fiber.Ctx) error {
	const op = "UploadCandidateResumeFileHandler"

	// auth := middleware.GetAuthData(c)
	// if auth == nil {
	// 	return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	// }

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrBadRequest.WithWrap(err).WithHints("form field `file` is required"),
			"%s.%s", ctrl.pkg, op,
		)
	}

	f, err := fileHeader.Open()
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrInternal.WithWrap(err),
			"%s.%s", ctrl.pkg, op,
		)
	}
	// nolint
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrInternal.WithWrap(err),
			"%s.%s", ctrl.pkg, op,
		)
	}

	resp, err := ctrl.crmTenantService.UploadCandidateResumeFile(c.Context(), &desc.UploadCandidateResumeFileRequest{
		File:     fileData,
		Filename: fileHeader.Filename,
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	fmt.Println(resp.FileId)

	// files, err := ctrl.tenantFacade.Account.UploadProfileImageFile(c.Context(), fileHeader.Filename, fileData)
	// if err != nil {
	// 	return err
	// }

	result := dto.UploadedFilePackDTO{
		// Files: make(map[string]*dto.FileDTO, len(files)),
	}

	// for _, file := range files {
	// 	result.Files[file.Target] = dto.NewFileDTO(file, ctrl.fileUC, true)
	// }

	return c.Status(fiber.StatusCreated).JSON(result)
}
