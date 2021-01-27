package normalizer

import "commerce/models"

// InitNormalizer will initialize normalizer helper
func InitNormalizer() Normalizer {
	return &normalizer{}
}

// Normalizer implements normalizing of data functions interface
type Normalizer interface {
	User(*models.User) interface{}
	Account(*models.Accounts) interface{}
	Balance(*models.UserBalance) interface{}
	Role(role *models.UserRole) interface{}
}

type normalizer struct{}

func (n *normalizer) User(user *models.User) interface{} {
	uRes := make(map[string]interface{})

	uRes["id"] = user.ID
	uRes["username"] = user.Username
	uRes["createdAt"] = user.CreatedAt
	uRes["account"] = nil
	uRes["balance"] = nil
	uRes["role"] = nil

	if user.Account.ID != 0 {
		uRes["account"] = n.Account(&user.Account)
	}

	if user.Balance.ID != 0 {
		uRes["balance"] = n.Balance(&user.Balance)
	}

	if user.Role.ID != 0 {
		uRes["role"] = n.Role(&user.Role)
	}

	return uRes
}

func (n *normalizer) Account(account *models.Accounts) interface{} {
	acc := make(map[string]interface{})

	acc["id"] = account.ID
	acc["firstName"] = account.FirstName
	acc["lastName"] = account.LastName
	acc["updatedOn"] = account.UpdatedAt

	return acc
}

func (n *normalizer) Balance(balance *models.UserBalance) interface{} {
	bal := make(map[string]interface{})

	bal["id"] = balance.ID
	bal["balance"] = balance.Balance
	bal["updatedOn"] = balance.UpdatedAt

	return bal
}

func (n *normalizer) Role(role *models.UserRole) interface{} {
	r := make(map[string]interface{})

	r["id"] = role.ID
	r["type"] = role.Type

	return r
}
