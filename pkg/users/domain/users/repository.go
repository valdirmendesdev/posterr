package users

type Repository interface {
	GetByUsername(un Username) (*User, error)
}
