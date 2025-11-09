package user

import (
	"context"
	"log/slog"

	"github.com/neurochar/backend/internal/app"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/app/fxboot/invoking"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

// Init - init domain
func Init(_ app.ID, logger *slog.Logger, roleUC userUC.RoleUsecase) invoking.InvokeInit {
	const op = "User.init"

	return invoking.InvokeInit{
		StartBeforeOpen: func(ctx context.Context) error {
			logger.InfoContext(ctx, "started to init domain", slog.String("domain", "user"))

			err := roleUC.BuildRolesInMemory(ctx)
			if err != nil {
				return appErrors.Chainf(err, "%s", op)
			}

			return nil
		},
	}
}
