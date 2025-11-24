package testing

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app"
	"github.com/neurochar/backend/internal/app/fxboot/invoking"
)

// Init - init domain
func Init(_ app.ID, logger *slog.Logger) invoking.InvokeInit {
	// const op = "Tenant.init"

	return invoking.InvokeInit{}
}
