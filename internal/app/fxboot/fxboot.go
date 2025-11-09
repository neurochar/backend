// Package fxboot contains fx bootstrapping
package fxboot

import (
	"go.uber.org/fx"
)

// ProvidingID - type for providing id
type ProvidingID int

const (
	// ProvidingAppID - app id
	ProvidingAppID ProvidingID = iota

	// ProvidingIDFXTimeouts - fx timeouts
	ProvidingIDFXTimeouts

	// ProvidingIDConfig - app config
	ProvidingIDConfig

	// ProvidingIDLogger - logger
	ProvidingIDLogger

	// ProvidingIDFXLogger - fx logger
	ProvidingIDFXLogger

	// ProvidingIDImageProc - image processor
	ProvidingIDImageProc

	// ProvidingIDDBClients - db clients
	ProvidingIDDBClients

	// ProvidingIDBackoff - backoff
	ProvidingIDBackoff

	// ProvidingIDStorageClient - storage client
	ProvidingIDStorageClient

	// ProvidingIDEmailing - emailing
	ProvidingIDEmailing

	// ProvidingIDJobsController - jobs controller
	ProvidingIDJobsController

	// ProvidingHTTPFiberServer - http fiber server
	ProvidingHTTPFiberServer

	// ProvidingIDDeliveryHTTP - delivery http
	ProvidingIDDeliveryHTTP

	// ProvidingIDFileModule - file module
	ProvidingIDFileModule

	// ProvidingIDUserModule - user module
	ProvidingIDUserModule

	// ProvidingIDEmailingModule - emailing module
	ProvidingIDEmailingModule

	// ProvidingIDAlertModule - alert
	ProvidingIDAlertModule

	// ProvidingIDTenantModule - tenant
	ProvidingIDTenantModule

	// ProvidingIDTenantUserModule - tenant user
	ProvidingIDTenantUserModule
)

// OptionsMap - options map item with providing and invokes elements
type OptionsMap struct {
	Providing map[ProvidingID]fx.Option
	Invokes   []fx.Option
}

// OptionsMapToSlice - convert options map to slice for fx bootstrapping
func OptionsMapToSlice(optionsMap OptionsMap) []fx.Option {
	result := make([]fx.Option, 0)

	for _, option := range optionsMap.Providing {
		result = append(result, option)
	}

	result = append(result, optionsMap.Invokes...)

	return result
}
