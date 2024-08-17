package services

import (
	"aswadwk/chatai/models"
	"aswadwk/chatai/repositories"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindUserBy(by, value string) (models.User, error)
	SaveUser(user models.User) error
	DeleteUser(user models.User) error
	GetUsers(query models.QueryModel) (interface{}, error)
	UpdateUser(user models.User) error
}

type userService struct {
	repo repositories.UserRepository
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(user models.User) error {
	if err := u.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

// GetUsers implements UserService.
func (u *userService) GetUsers(query models.QueryModel) (interface{}, error) {
	response, err := u.repo.GetUsersWithPaginate(query)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteUser implements UserService.
func (u *userService) DeleteUser(user models.User) error {
	err := u.repo.DeleteUser(user)

	if err != nil {
		return err
	}

	return nil
}

// SaveUser implements UserService.
func (u *userService) SaveUser(user models.User) error {

	hash, err := hashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hash

	err = u.repo.SaveUser(user)

	if err != nil {
		return err
	}

	return nil
}

// FindUserBy implements UserService.
func (u *userService) FindUserBy(by string, value string) (models.User, error) {
	user, err := u.repo.FindUserBy(by, value)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return string(hash), nil
}
