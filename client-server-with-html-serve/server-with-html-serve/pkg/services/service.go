package service

import (
	repo "github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/repositories"
)

func NewService(repo *repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo *repo.Repo
}
