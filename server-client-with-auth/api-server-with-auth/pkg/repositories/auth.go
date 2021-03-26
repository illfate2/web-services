package repo

import (
	"github.com/jackc/pgx/v4"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

type Auth struct {
	ID           int
	ProviderID   string `xorm:"provider_id"`
	ProviderType string `xorm:"provider"`
	UserID       int    `xorm:"user_id"`
}

func (a Auth) toEntity() entities.Auth {
	return entities.Auth{
		ID:             a.ID,
		ProviderUserID: a.ProviderID,
		ProviderType:   a.ProviderType,
		UserID:         a.UserID,
	}
}

func (r *Repo) InsertAuth(auth entities.Auth) (entities.Auth, error) {
	repoAuth := Auth{
		ProviderID:   auth.ProviderUserID,
		ProviderType: auth.ProviderType,
		UserID:       auth.UserID,
	}
	id, err := r.engine.InsertOne(repoAuth)
	repoAuth.ID = int(id)
	return repoAuth.toEntity(), err
}

func (r *Repo) FindAuthByProvider(providerID, providerType string) (entities.Auth, error) {
	var auth Auth
	ok, err := r.engine.Where("provider_id = ? AND provider = ?", providerID, providerType).Get(&auth)
	if err != nil {
		return entities.Auth{}, err
	}
	if !ok {
		return entities.Auth{}, pgx.ErrNoRows
	}
	return entities.Auth{
		ID:             auth.ID,
		ProviderUserID: auth.ProviderID,
		ProviderType:   auth.ProviderType,
		UserID:         auth.UserID,
	}, nil
}
