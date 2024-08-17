package services

import (
	"time"

	"aswadwk/chatai/config"
	"aswadwk/chatai/models"
	"aswadwk/chatai/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateToken(user models.User) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	repo repositories.UserRepository
}

// GenerateToken implements JwtService.
func (j *jwtService) GenerateToken(user models.User) (string, error) {
	user, err := j.repo.FindUserBy("email", user.Email)

	if err != nil {
		return "", err
	}

	token, err := generateToken(user)

	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken implements JwtService.
func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	decode, err := parseToken(token)

	if err != nil {
		return nil, err
	}

	return decode, nil
}

func NewJwtService(repo repositories.UserRepository) JwtService {
	return &jwtService{
		repo: repo,
	}
}

func generateToken(user models.User) (string, error) {
	ttl := config.Config("JWT_ACCESS_TOKEN_TTL")

	if ttl == "" {
		panic("JWT_ACCESS_TOKEN_TTL is not set")
	}
	secretKey := config.Config("JWT_SECRET_KEY_ACCESS_TOKEN")

	// convert ttl int
	ttlDuration, err := time.ParseDuration(ttl)

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iss":  "x-chat",
		"sub":  user.Email,
		"aud":  "x-chat",
		"exp":  time.Now().Add(ttlDuration).Unix(),
		"nbf":  time.Now().Unix(),
		"iat":  time.Now().Unix(),
		"jti":  user.ID,
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func parseToken(token string) (*jwt.Token, error) {
	secretKey := config.Config("JWT_SECRET_KEY_ACCESS_TOKEN")

	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Token invalid or expired")
	}

	return tokenParse, nil
}
