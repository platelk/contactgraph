package userquery

import (
	"github.com/platelk/contactgraph/domain/models/users"
)

// Queryer will construct the query to search users based on multiple criteria
type Queryer interface {
	ByID(id users.ID) Queryer
	ByPhoneNumber(phoneNumber users.PhoneNumber) Queryer
	ByNickName(nickName string) Queryer
}
