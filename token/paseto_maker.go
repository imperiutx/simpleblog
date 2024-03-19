package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoMaker(symmetricKey paseto.V4SymmetricKey, implicit string) (*PasetoMaker, error) {
	return &PasetoMaker{
		symmetricKey: symmetricKey,
		implicit:     []byte(implicit),
	}, nil
}
func (maker *PasetoMaker) CreateToken(username, role string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", err
	}

	token := paseto.NewToken()
	token.Set("id", payload.ID.String())
	token.Set("username", payload.Username)
	token.Set("role", payload.Role)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	// Encrypt and return the token
	encryptedToken := token.V4Encrypt(maker.symmetricKey, maker.implicit)

	return encryptedToken, nil
}

func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	parser := paseto.NewParser()
	rules := []paseto.Rule{
		paseto.NotExpired(),
	}
	parser.SetRules(rules)

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, tokenString, maker.implicit)
	if err != nil {
		return nil, err
	}

	payload, err := getPayload(parsedToken)
	if err != nil {
		return nil, err
	}

	return payload, nil

}

func getPayload(token *paseto.Token) (*Payload, error) {
	payload := &Payload{}

	var err error

	tokenid, err := token.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload.ID, err = uuid.Parse(tokenid)
	if err != nil {
		return nil, ErrInvalidToken
	}
	payload.Username, err = token.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}
	payload.IssuedAt, err = token.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	payload.ExpiredAt, err = token.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
