package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

// CreateToken implements Maker
func (pasetoMaker *PasetoMaker) CreateToken(email string, duration time.Duration) (string, *Payload, error) {
	payload,err:=NewPayload(email,duration);
	if err!=nil{
		return "",payload,err
	}
	token,err:=pasetoMaker.paseto.Encrypt(pasetoMaker.symetricKey,payload,nil)
	return token,payload,err
}

// VerifyToken implements Maker
func (pasetoMaker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload:=&Payload{}
	err:=pasetoMaker.paseto.Decrypt(token,pasetoMaker.symetricKey,payload,nil)
	if (err!=nil){
		return nil,ErrInvalidToken
	}
	err=payload.Valid();
	if err!=nil{
		return nil,err
	}
	return payload,nil
}

func NewPasetoMaker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be exactly %d length", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(key),
	}, nil
}
