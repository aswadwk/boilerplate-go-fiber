package controllers

import (
	"fmt"
	"strconv"

	userDto "aswadwk/chatai/dto/user"
	"aswadwk/chatai/helpers"
	"aswadwk/chatai/models"
	"aswadwk/chatai/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	FindUserBy(ctx *fiber.Ctx) error
	SaveUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
	GetUsers(ctx *fiber.Ctx) error
}

type userController struct {
	service  services.UserService
	validate *validator.Validate
}

// GetUsers implements UserController.
func (u *userController) GetUsers(ctx *fiber.Ctx) error {
	tenantId := ctx.Locals("current_tenant_id").(string)
	page, err := strconv.Atoi(ctx.Query("page", "1"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	query := models.QueryModel{
		All:      false,
		Page:     page,
		PerPage:  perPage,
		TenantID: tenantId,
	}

	users, err := u.service.GetUsers(query)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.Success(
		"Success", users, nil,
	))
}

// DeleteUser implements UserController.
func (u *userController) DeleteUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	currentUserId := ctx.Locals("current_user_id").(string)

	// check role must be super admin or admin
	currentRole := ctx.Locals("current_role").(string)

	if currentRole != "super admin" && currentRole != "admin" {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(helpers.Error("Forbidden", nil, nil))
	}

	if currentUserId == userId {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(helpers.Error("Cannot delete admin", nil, nil))
	}

	err := u.service.DeleteUser(models.User{
		ID: userId,
	})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(helpers.Error("Failed to delete user", err, nil))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(helpers.Success("Success to delete user", nil, nil))
}

// SaveUser implements UserController.
func (u *userController) SaveUser(ctx *fiber.Ctx) error {
	newUser := userDto.NewUserDto{}

	if err := ctx.BodyParser(&newUser); err != nil {
		return ctx.Status(400).
			JSON(helpers.Error("Failed to parse request body", err, nil))
	}

	currentRole := ctx.Locals("current_role").(string)

	if currentRole != "super admin" {
		newUser.Role = "user"
		newUser.TenantID = ctx.Locals("current_tenant_id").(string)
	}

	if currentRole == "user" {
		return ctx.Status(403).JSON(helpers.Error("Forbidden", nil, nil))
	}

	fmt.Println("new user", newUser)

	err := u.validate.Struct(newUser)

	if err != nil {
		fmt.Println("error", err)

		return ctx.Status(400).
			JSON(helpers.Error(
				"Failed to validate request body", err, nil,
			))
	}

	err = u.service.SaveUser(models.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		Role:     newUser.Role,
	})

	if err != nil {
		return ctx.Status(400).
			JSON(helpers.Error("Failed to save user", err, nil))
	}

	return ctx.Status(200).
		JSON(helpers.Success("Success to save user", nil, nil))
}

// FindUserBy implements UserController.
func (u *userController) FindUserBy(ctx *fiber.Ctx) error {
	by := ctx.Params("by")
	value := ctx.Params("value")

	user, err := u.service.FindUserBy(by, value)

	if err != nil {
		ctx.Status(500)
		ctx.JSON(fiber.Map{
			"message": "Failed to find user",
		})

		return err
	}

	ctx.Status(200)
	ctx.JSON(fiber.Map{
		"message": "Success to find user",
		"data":    user,
	})

	return nil
}

func NewUserController(service services.UserService, validate *validator.Validate) UserController {
	return &userController{
		service:  service,
		validate: validate,
	}
}
