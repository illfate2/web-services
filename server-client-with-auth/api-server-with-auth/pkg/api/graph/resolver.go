package graph

import (
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/auth"
	service "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service    *service.Service
	jwtService *auth.JWTService
}

func NewResolver(service *service.Service, jwtService *auth.JWTService) *Resolver {
	return &Resolver{service: service, jwtService: jwtService}
}
