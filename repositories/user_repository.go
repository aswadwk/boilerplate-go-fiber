package repositories

import (
	"fmt"

	"aswadwk/chatai/helpers"
	"aswadwk/chatai/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user models.User) error
	DeleteUser(user models.User) error
	FindUserBy(by, value string) (models.User, error)
	GetUsersWithPaginate(query models.QueryModel) (models.QueryModelResponse, error)
	UpdateUser(user models.User) error
}

type userRepository struct {
	db *gorm.DB
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(user models.User) error {
	err := u.db.Model(&user).Updates(&user).Error

	if err != nil {
		return err
	}

	return nil
}

// GetUsersWithPaginate implements UserRepository.
func (u *userRepository) GetUsersWithPaginate(query models.QueryModel) (models.QueryModelResponse, error) {
	users := []models.User{}

	response, err := helpers.QueryPaginate(
		u.db,
		&models.User{},
		users,
		query.Page,
		query.PerPage,
		helpers.SearchBy("tenant_id", query.TenantID),
	)

	if err != nil {
		return response, err
	}

	return response, nil
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(user models.User) error {
	err := u.db.Delete(&user).Error

	if err != nil {
		return err
	}

	return nil
}

// FindUserBy implements UserRepository.
func (u *userRepository) FindUserBy(by, value string) (models.User, error) {
	user := models.User{}

	if err := u.db.Where(fmt.Sprintf("%s = ?", by), value).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// SaveUser implements UserRepository.
func (u *userRepository) SaveUser(user models.User) error {

	user.ID = uuid.New().String()

	err := u.db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
