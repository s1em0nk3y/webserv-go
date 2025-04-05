// Интерфейсы для работы приложения
package app

// Интерфейс для работы с пользователями
type UserStorage interface {
	CreateUser(user, passwordHash string) error
}
