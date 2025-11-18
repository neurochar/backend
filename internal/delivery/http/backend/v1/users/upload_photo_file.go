package users

import (
	"io"

	"github.com/gofiber/fiber/v2"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/dto"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
)

func (ctrl *Controller) UploadPhotoFileHandler(c *fiber.Ctx) error {
	const op = "UploadFileHandler"

	auth := middleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	isRevoked, err := ctrl.tenantUserFacade.Auth.IsSessionRevoked(c.Context(), auth.SessionID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if isRevoked {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

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

	files, err := ctrl.tenantUserFacade.Account.UploadProfileImageFile(c.Context(), fileHeader.Filename, fileData)
	if err != nil {
		return err
	}

	result := dto.UploadedFilePackDTO{
		Files: make(map[string]*dto.FileDTO, len(files)),
	}

	for _, file := range files {
		result.Files[file.Target] = dto.NewFileDTO(file, ctrl.fileUC, true)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
