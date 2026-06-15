package services

import (
	"fmt"
	"rompelago/auth"
	"rompelago/dto"
	"rompelago/repositories"
	"strconv"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepositories
}

func InitAuthService(userRepo *repositories.UserRepositories) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func isValidPassword(p string) bool {
	if len(p) < 12 {
		return false
	}
	var hasUpper, hasSpecial bool
	for _, c := range p {
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			hasSpecial = true
		}
	}
	return hasUpper && hasSpecial
}

// Login vérifie les credentials et retourne un LoginResponseDto avec le JWT.
func (s *AuthService) Login(req dto.LoginRequestDto) (dto.LoginResponseDto, error) {
	if req.Pseudo == "" || req.Password == "" {
		return dto.LoginResponseDto{}, fmt.Errorf("pseudo et mot de passe obligatoires")
	}

	user, hashedPassword, err := s.userRepo.FindByPseudo(req.Pseudo)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("identifiants incorrects")
	}

	// La ligne en dessous est pas si compliquée que ça
	// Ce que je fais c'est que je récupère le mdp depuis la DB en hashé
	// Je hash le mdp que l'utilisateur a rentré et je compare les deux
	// Pas si compliqué que ça mais je sais que Nicolas comprendras pas donc j'explique
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("identifiants incorrects")
	}

	if user.StatusBan == 1 {
		return dto.LoginResponseDto{}, fmt.Errorf("compte banni")
	}

	token, err := auth.GenerateToken(strconv.Itoa(user.Id), user.Role)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("erreur génération token - %v", err)
	}

	return dto.LoginResponseDto{
		Type:        "Bearer",
		AccessToken: token,
		ExpiresIn:   1500,
	}, nil
}

// Register crée un nouvel utilisateur après validation.
func (s *AuthService) Register(req dto.RegisterRequestDto) error {
	if req.Pseudo == "" || req.Email == "" || req.Password == "" {
		return fmt.Errorf("pseudo, email et mot de passe obligatoires")
	}
	if !isValidPassword(req.Password) {
		return fmt.Errorf("le mot de passe doit contenir au moins 12 caractères, dont une majuscule et un caractère spécial")
	}

	exists, err := s.userRepo.ExistsByPseudoOrEmail(req.Pseudo, req.Email)
	if err != nil {
		return fmt.Errorf("erreur vérification - %v", err)
	}
	if exists {
		return fmt.Errorf("pseudo ou email déjà utilisé")
	}

	// Hashage bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("erreur hashage mot de passe - %v", err)
	}

	_, err = s.userRepo.Create(req.Pseudo, req.Email, string(hashed))
	return err
}
