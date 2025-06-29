package webauthn

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/domain/user"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*user.User, error)
	Put(ctx context.Context, u *user.User) (*user.User, error)
}

type Service struct {
	wa       *webauthn.WebAuthn
	userRepo UserRepository
}

func New(cfg *config.Config, userRepo UserRepository) *Service {
	proto := cfg.HTTPServer.Protocol
	host := cfg.HTTPServer.Host
	port := cfg.HTTPServer.Port
	origin := fmt.Sprintf("%s://%s", proto, host)
	originWithPort := fmt.Sprintf("%s://%s%s", proto, host, port)

	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Port Service",
		RPID:          host,
		RPOrigins:     []string{origin, originWithPort},
	})

	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		wa:       wa,
		userRepo: userRepo,
	}
}

func (s *Service) BeginRegistration(ctx context.Context, id string) (
	creation *protocol.CredentialCreation, session *webauthn.SessionData, err error) {
	var u *user.User
	u, err = s.userRepo.Get(ctx, id) // Find or create the new user
	if err != nil {
		u, err = user.New(id, id, id)
		if err != nil {
			return nil, nil, err
		}
		_, err = s.userRepo.Put(ctx, u) // Save the new user
		if err != nil {
			return nil, nil, err
		}
	}

	creation, session, err = s.wa.BeginRegistration(u)
	if err != nil {
		return nil, nil, err
	}

	return
}

func (s *Service) FinishRegistration(
	ctx context.Context,
	session webauthn.SessionData,
	r *http.Request,
) error {
	id := string(session.UserID)
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	credential, err := s.wa.FinishRegistration(user, session, r)
	if err != nil {
		return fmt.Errorf("failed to finish registration: %w", err)
	}

	user.AddCredential(credential)
	_, err = s.userRepo.Put(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to save user credentials: %w", err)
	}

	return nil
}

func (s *Service) BeginLogin(
	ctx context.Context,
	id string,
) (creation *protocol.CredentialAssertion, session *webauthn.SessionData, err error) {
	var u *user.User
	u, err = s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	creation, session, err = s.wa.BeginLogin(u)
	if err != nil {
		return nil, nil, err
	}

	return
}

func (s *Service) FinishLogin(
	ctx context.Context,
	session webauthn.SessionData,
	r *http.Request,
) error {
	user, err := s.userRepo.Get(ctx, string(session.UserID))
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	credential, err := s.wa.FinishLogin(user, session, r)
	if err != nil {
		return fmt.Errorf("failed to finish registration: %w", err)
	}

	if credential.Authenticator.CloneWarning {
		// TODO: Handle clone warning
		log.Println("Authenticator clone warning detected")
	}

	user.AddCredential(credential)
	_, err = s.userRepo.Put(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to save user credentials: %w", err)
	}

	return nil
}
