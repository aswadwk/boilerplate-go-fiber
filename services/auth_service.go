package services

import (
	"fmt"

	"aswadwk/chatai/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, error)
	CurrentUser(token string) (models.User, error)
	ChangePassword(user models.User, oldPassword, newPassword string) error
}

type authService struct {
	userService UserService
	jwtService  JwtService
}

// ChangePassword implements AuthService.
func (a *authService) ChangePassword(user models.User, oldPassword string, newPassword string) error {
	checkUser := models.User{}

	checkUser, err := a.userService.FindUserBy("email", user.Email)
	if err != nil {
		return err
	}

	if !verifyPassword(checkUser.Password, oldPassword) {
		return fmt.Errorf("invalid credentials")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err = a.userService.UpdateUser(models.User{
		ID:       user.ID,
		Password: string(hashedPassword),
	}); err != nil {
		return err
	}

	return nil
}

// CurrentUser implements AuthService.
func (a *authService) CurrentUser(token string) (models.User, error) {
	user := models.User{}

	decode, err := a.jwtService.ValidateToken(token)

	if err != nil {
		return models.User{}, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	issuer, err := decode.Claims.GetSubject()

	if err != nil {
		return user, err
	}

	user, err = a.userService.FindUserBy("email", issuer)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Login implements AuthService.
func (a *authService) Login(email string, password string) (string, error) {
	user, err := a.userService.FindUserBy("email", email)

	if err != nil {
		return "", err
	}

	if !verifyPassword(user.Password, password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := a.jwtService.GenerateToken(user)

	if err != nil {
		return "", err
	}

	return token, nil
}

func NewAuthService(userService UserService, jwtService JwtService) AuthService {
	return &authService{
		userService: userService,
		jwtService:  jwtService,
	}
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
