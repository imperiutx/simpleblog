package token

import (
	"encoding/json"
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(hexKey string) (*PasetoMaker, error) {
	symmetric, err := paseto.V4SymmetricKeyFromHex(hexKey)
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		symmetricKey: symmetric,
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}

	jsonClaim, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}

	token, err := paseto.NewTokenFromClaimsJSON(jsonClaim, []byte{})
	if err != nil {
		return "", nil, err
	}
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	// key := paseto.NewV4SymmetricKey() // don't share this!!
	encrypted := token.V4Encrypt(maker.symmetricKey, nil)
	return encrypted, payload, nil
}

func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	parser := paseto.NewParser()
	rules := []paseto.Rule{
		paseto.NotExpired(),
	}
	parser.SetRules(rules)

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, err
	}

	payload := &Payload{}
	if err := json.Unmarshal(parsedToken.ClaimsJSON(), payload); err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil

}
