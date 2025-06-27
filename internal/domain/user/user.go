package user

import (
	"errors"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
)

// TODO: extract to a separate package
var (
	ErrNotFound   = errors.New("user not found")
	ErrValidation = errors.New("validation error")
	ErrRequired   = fmt.Errorf("%w: value cannot be empty", ErrValidation)
)

type User struct {
	ID          []byte
	DisplayName string
	Name        string

	creds []webauthn.Credential
}

func New(id string, name, displayName string) (*User, error) {
	if id == "" {
		return nil, ErrRequired
	}
	if name == "" {
		return nil, ErrRequired
	}
	if displayName == "" {
		return nil, ErrRequired
	}

	return &User{
		ID:          []byte(id),
		Name:        name,
		DisplayName: displayName,
	}, nil
}

func (o *User) WebAuthnID() []byte                            { return o.ID }
func (o *User) WebAuthnName() string                          { return o.Name }
func (o *User) WebAuthnDisplayName() string                   { return o.DisplayName }
func (o *User) WebAuthnCredentials() []webauthn.Credential    { return o.creds }
func (o *User) AddCredential(credential *webauthn.Credential) { o.creds = append(o.creds, *credential) }
func (o *User) UpdateCredential(credential *webauthn.Credential) {
	for i, c := range o.creds {
		if string(c.ID) == string(credential.ID) {
			o.creds[i] = *credential
		}
	}
}
