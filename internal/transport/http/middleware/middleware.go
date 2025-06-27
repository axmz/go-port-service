package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/alexedwards/scs/v2"
	wah "github.com/axmz/go-port-service/internal/transport/http/handlers/webauthn"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	RequestIDHeader = "X-Request-Id"
	reqid           uint64
)

type requestIdCtxType int

const (
	RequestIDKey requestIdCtxType = 0
)

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			id := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%06d", id)
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func Logger(next http.Handler) http.Handler {
	const op = "transport.http.middleware.Logger"
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := GetReqID(r.Context())

		slog.With(
			slog.String("component", op),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote", r.RemoteAddr),
			slog.String("agent", r.UserAgent()),
			slog.String("req_id", reqID),
		).Info("started request")

		defer func() {
			slog.With(
				slog.String("req_id", reqID),
				slog.String("duration", time.Since(start).String()),
			).Info("completed request")
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Recoverer(next http.Handler) http.Handler {
	const op = "transport.http.middleware.Recoverer"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := GetReqID(r.Context())

				slog.Error("panic recovered",
					slog.String("op", op),
					slog.Any("error", rec),
					slog.String("stack", string(debug.Stack())),
					slog.String("req_id", reqID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
				)

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func LoggedInMiddleware(session *scs.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := session.Get(r.Context(), wah.WebauthSessionKey).(webauthn.SessionData)
		if !ok {
			slog.Error("session not found",
				slog.String("op", "LoggedInMiddleware"),
				slog.String("req_id", GetReqID(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if s.UserID == nil { // check expired?
			// Option 1: return 401 Unauthorized
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// Option 2: redirect to login page
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
