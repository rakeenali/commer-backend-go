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
	Tag(*models.Tags, bool) interface{}
	Item(*models.Items) interface{}
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

func (n *normalizer) Tag(tag *models.Tags, withItem bool) interface{} {
	r := make(map[string]interface{})

	r["id"] = tag.ID
	r["name"] = tag.Name
	r["updateAt"] = tag.UpdatedAt
	r["createdAt"] = tag.CreatedAt
	r["items"] = nil
	r["totalItems"] = len(tag.Items)

	if len(tag.Items) > 0 && withItem {
		var items []interface{}
		for _, i := range tag.Items {
			item := n.Item(&i)
			items = append(items, item)
		}
		r["items"] = items
	}

	return r
}

func (n *normalizer) Item(item *models.Items) interface{} {
	r := make(map[string]interface{})

	r["id"] = item.ID
	r["name"] = item.Name
	r["image"] = item.Image
	r["price"] = item.Price
	r["sku"] = item.Sku
	r["createAt"] = item.CreatedAt
	r["updatedAt"] = item.UpdatedAt
	r["tags"] = nil

	if len(item.Tags) > 0 {
		var tags []interface{}
		for _, t := range item.Tags {
			tag := n.Tag(&t, true)
			tags = append(tags, tag)
		}
		r["tags"] = tags
	}

	return r
}
