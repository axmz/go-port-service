package webauthn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/axmz/go-port-service/internal/transport/http/response"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

const WebauthSessionKey = "webauthn_session"

type WebAuthnService interface {
	BeginRegistration(ctx context.Context, userID string) (*protocol.CredentialCreation, *webauthn.SessionData, error)
	FinishRegistration(ctx context.Context, session webauthn.SessionData, r *http.Request) error
	BeginLogin(ctx context.Context, userID string) (*protocol.CredentialAssertion, *webauthn.SessionData, error)
	FinishLogin(ctx context.Context, session webauthn.SessionData, r *http.Request) error
}

type SessionManager interface {
	Put(ctx context.Context, key string, data any)
	Get(ctx context.Context, key string) any
}

type Handlers struct {
	webauthn WebAuthnService
	session  SessionManager
}

func New(webauthn WebAuthnService, session SessionManager) *Handlers {
	return &Handlers{
		webauthn: webauthn,
		session:  session,
	}
}

func getUserID(r *http.Request) (string, error) {
	type UserID struct {
		UserID string `json:"email"`
	}
	var u UserID
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return "", err
	}

	return u.UserID, nil
}

func (h *Handlers) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		panic(err) // FIXME: handle error
	}

	options, session, err := h.webauthn.BeginRegistration(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("can't begin registration: %s", err.Error())
		response.BadRequest(w, msg)
		return
	}

	// Make a session key and store the sessionData values
	h.session.Put(r.Context(), WebauthSessionKey, session)

	// return the options generated with the session key
	response.OK(w, options)
	// options.publicKey contain our registration options
}

func (h *Handlers) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	session := h.session.Get(r.Context(), WebauthSessionKey).(webauthn.SessionData)

	err := h.webauthn.FinishRegistration(r.Context(), session, r)
	if err != nil {
		msg := fmt.Sprintf("can't finish registration: %s", err.Error())
		response.BadRequest(w, msg)
		return
	}

	response.OK(w, "Registration Success")
}

func (h *Handlers) BeginLogin(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		panic(err) // FIXME: handle error
	}

	options, session, err := h.webauthn.BeginLogin(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("can't begin login: %s", err.Error())
		response.BadRequest(w, msg)
		return
	}

	// Make a session key and store the sessionData values
	h.session.Put(r.Context(), WebauthSessionKey, session)

	// return the options generated with the session key
	response.OK(w, options)
	// options.publicKey contain our registration options

}
func (h *Handlers) FinishLogin(w http.ResponseWriter, r *http.Request) {
	session := h.session.Get(r.Context(), WebauthSessionKey).(webauthn.SessionData)

	err := h.webauthn.FinishLogin(r.Context(), session, r)
	if err != nil {
		msg := fmt.Sprintf("can't finish registration: %s", err.Error())
		response.BadRequest(w, msg)
		return
	}

	response.OK(w, "Registration Success")
}
