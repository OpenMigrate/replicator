package middleware

import (
	"context"
	"net/http"
	"replicator/internal/storage"
)

func WithStore(s *storage.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), storeKey, s)))
		})
	}
}

type ctxKey string

const storeKey ctxKey = "store"

func StoreFrom(r *http.Request) *storage.Store {
	v := r.Context().Value(storeKey)
	if v == nil {
		return nil
	}
	return v.(*storage.Store)
}
