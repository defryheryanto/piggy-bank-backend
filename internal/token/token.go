package token

type TokenIService interface {
	GenerateToken(payload interface{}) (string, error)
	Parse(token string) (interface{}, error)
	CheckValidity(token string) (bool, error)
}
