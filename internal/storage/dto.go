// storage/dto.go
package storage

type AppCreate struct {
	ID          string
	Name        string
	Identifier  string
	Description string
}

type AppSelector struct {
	ID         *string
	Identifier *string
}

type Cursor struct {
	AfterID string
	Limit   int
}

type MembershipOp string

const (
	MembershipAdd     MembershipOp = "add"
	MembershipRemove  MembershipOp = "remove"
	MembershipReplace MembershipOp = "replace"
)
