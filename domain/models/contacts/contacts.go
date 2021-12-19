package contacts

import (
	"github.com/platelk/contactgraph/domain/models/users"
)

type void struct{}
type List map[users.ID]void

func (s List) ToSlice() []users.ID {
	var l []users.ID
	for k := range s {
		l = append(l, k)
	}
	return l
}
