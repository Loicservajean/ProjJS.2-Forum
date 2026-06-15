package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken verifie un JWT et retourne ses informations si le token est valide.
func ValidateToken(tokenString string) (*Claims, error) {
	// Claims recevra les donnees decodees depuis le token.
	claims := &Claims{}

	// ParseWithClaims analyse le token, verifie sa signature et remplit claims.
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// On refuse tout algorithme different de HS256 pour eviter les signatures inattendues.
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			// Cette cle doit etre la meme que celle utilisee pour signer le token.
			return getSecret(), nil
		},
		// Ces options garantissent que le token vient de l'API attendue
		// et qu'il est destine au bon client.
		jwt.WithIssuer("product-api"),
		jwt.WithAudience("product-front"),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	// Une erreur ici signifie que le token est mal forme, expire ou invalide.
	if err != nil {
		return nil, err
	}

	// Par securite, on verifie aussi le statut final du token parse.
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Le token est valide : on retourne les claims utilisables par l'application.
	return claims, nil
}
