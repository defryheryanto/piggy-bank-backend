package token

type TokenIService[T any] interface {
	GenerateToken(payload T) (string, error)
	Parse(token string) (T, error)
	CheckValidity(token string) (bool, error)
}
