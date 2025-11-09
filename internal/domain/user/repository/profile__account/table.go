package profileaccount

import (
	accountRepo "github.com/neurochar/backend/internal/domain/user/repository/account"
	profileRepo "github.com/neurochar/backend/internal/domain/user/repository/profile"
)

var JoinFields = []string{}

func init() {
	for _, field := range accountRepo.AccountTableFields {
		JoinFields = append(JoinFields, "a."+field+" as account___"+field)
	}

	for _, field := range profileRepo.ProfileTableFields {
		JoinFields = append(JoinFields, "p."+field+" as profile___"+field)
	}
}
