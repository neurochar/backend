package profileaccount

import (
	accountRepo "github.com/neurochar/backend/internal/domain/user/repository/account"
	profileRepo "github.com/neurochar/backend/internal/domain/user/repository/profile"
	"github.com/neurochar/backend/internal/domain/user/usecase"
)

// DBModel - database model
type DBModel struct {
	Profile profileRepo.DBModel `db:"profile"`
	Account accountRepo.DBModel `db:"account"`
}

func (db *DBModel) ToEntity() *usecase.User {
	return &usecase.User{
		Profile: db.Profile.ToEntity(),
		Account: db.Account.ToEntity(),
	}
}
