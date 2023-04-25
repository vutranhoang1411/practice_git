package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct{
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(email string,duration time.Duration)(*Payload,error){
	tokenID,err:=uuid.NewUUID()
	if (err!=nil){
		return nil,err
	}
	payload:=&Payload{
		ID:tokenID,
		Email: email,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload,nil
}
func (payload *Payload)Valid()error{
	if time.Now().After(payload.ExpiredAt){
		return ErrExpiredToken
	}
	return nil;
}