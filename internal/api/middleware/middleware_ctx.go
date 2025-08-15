package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"replicator/internal/storage"
)

type ctxKey string
const storeKey ctxKey = "store"

// Middleware func, updates db sotore key & it's reference in it's context
func WithStore(s *storage.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), storeKey, s)))
		})
	}
}

func InjectLog(log *slog.Logger) func(http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger" , log)))
    })
  }
}

func GetLogFromCtx(r *http.Request) (log *slog.Logger){
  logger := r.Context().Value("logger")

  // type check
  log, _ = logger.(*slog.Logger)
  return
}

func StoreFrom(r *http.Request) (s *storage.Store) {
	v := r.Context().Value(storeKey) // finding with the db with key "store"
	if v == nil {
		return nil
	}
  s, _ = v.(*storage.Store) // validating the storage type
	return
}
