package middleware

import (
	"context"
	"net/http"
	"rompelago/auth"
)

// UserContextKey est la clé utilisée pour stocker les claims dans le contexte de la requête.
type contextKey string

const UserContextKey contextKey = "user"

// WebAuthMiddleware lit le cookie "token", valide le JWT si présent et ajoute les
// claims au contexte. Contrairement à AuthMiddleware (API), il ne bloque jamais la
// requête : si le cookie est absent ou invalide, la requête continue sans claims
// (l'utilisateur est simplement considéré comme non connecté).
func WebAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			// Token invalide ou expiré : on continue sans utilisateur connecté.
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth protège une route web : redirige vers /connection si l'utilisateur
// n'est pas authentifié (pas de claims dans le contexte).
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(UserContextKey).(*auth.Claims); !ok {
			http.Redirect(w, r, "/connection", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
