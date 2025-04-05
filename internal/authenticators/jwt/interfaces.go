package jwt

type UserStorage interface {
	CreateUser(user, passwordHash string) error
	CheckCredents(user, passwordHash string) (bool, error)
	ValidateUsername(username string) error
}
