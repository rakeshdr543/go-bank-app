package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size")
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken implements Maker.
func (p *PasetoMaker) CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)

	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

// VerifyToken implements Maker.
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()

	if err != nil {
		return nil, err
	}

	return payload, nil

}
