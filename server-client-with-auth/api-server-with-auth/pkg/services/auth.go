package service

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

func (s *Service) SignupWithOAuth(providerID, providerType string) (entities.User, error) {
	provider, err := s.repo.FindAuthByProvider(providerID, providerType)
	if err == pgx.ErrNoRows {
		return s.createOAuthUser()
	}
	if err != nil {
		return entities.User{}, err
	}
	return s.repo.FindUser(provider.UserID)
}

func (s *Service) createOAuthUser() (entities.User, error) {
	return s.repo.InsertEmptyUser(entities.User{})
}

func (s *Service) CreateUser(user entities.User) (entities.User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, err
	}

	user.Password = string(pwdHash)
	u, err := s.repo.InsertUser(user)
	if err != nil {
		return entities.User{}, err
	}
	return u, nil
}

var ErrIncorrectPwd = errors.New("incorrect password")

func (s *Service) GetUser(email, password string) (entities.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return entities.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return entities.User{}, ErrIncorrectPwd
	}
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}
