package middleware

import (
	"context"
	"net/http"
	"rompelago/auth"
	"rompelago/helper"
	"strings"
)

// AuthMiddleware protege une route en exigeant un JWT valide.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Le token doit etre envoye dans le header Authorization.
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			helper.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Le format attendu est : Bearer <token>.
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helper.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		// On valide le token et on recupere les claims de l'utilisateur.
		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			helper.WriteError(w, http.StatusUnauthorized, "nvalid token")
			return
		}

		// Les claims sont ajoutes au contexte pour etre reutilises par les handlers suivants.
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
