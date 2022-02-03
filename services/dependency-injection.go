//go:build wireinject
// +build wireinject

package services

import (
	"sync"

	"github.com/google/wire"
)

type ServiceContainer struct {
	DB        IDbOperations
	Shortener IShortener
}

var containerOnce sync.Once
var Container *ServiceContainer

// @ref: https://github.com/google/wire/issues/260#issuecomment-680134730
func dependencyInjection() *ServiceContainer {
	wire.Build(newDB, dBConn, newShortener, wire.Struct(new(ServiceContainer), "*"))
	return nil
}

func RegisterServiceContainer() *ServiceContainer {
	containerOnce.Do(func() {
		Container = dependencyInjection()
	})

	return Container
}
