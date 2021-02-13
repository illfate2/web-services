package graph

import service "github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service *service.Service
}

func NewResolver(s *service.Service) *Resolver {
	return &Resolver{
		service: s,
	}
}
