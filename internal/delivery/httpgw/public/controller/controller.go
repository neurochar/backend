package controller

import (
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/auth_tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/crm"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/registration"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/rooms"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/testing"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/users_tenant"
)

const maxUploadSize = 50 << 20

type Controls struct {
	AuthTenant   *auth_tenant.Controller
	CRM          *crm.Controller
	Tenant       *tenant.Controller
	Registration *registration.Controller
	Rooms        *rooms.Controller
	UsersTenant  *users_tenant.Controller
	Testing      *testing.Controller
}

type Controller struct {
	pkg      string
	controls *Controls
}

func New(controls *Controls) *Controller {
	return &Controller{
		pkg:      "httpGateway.Controller",
		controls: controls,
	}
}

func (ctrl *Controller) Controls() *Controls {
	return ctrl.controls
}
