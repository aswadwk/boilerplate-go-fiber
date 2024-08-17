package controllers

import (
	"strings"

	userDto "aswadwk/chatai/dto/user"
	"aswadwk/chatai/helpers"
	"aswadwk/chatai/models"
	"aswadwk/chatai/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	CurrentUser(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
}

type authController struct {
	service  services.AuthService
	validate *validator.Validate
}

// ChangePassword implements AuthController.
func (a *authController) ChangePassword(ctx *fiber.Ctx) error {
	changePasswordDto := userDto.ChangePasswordDto{}

	user := models.User{}

	user.ID = ctx.Locals(helpers.CurrentUserID).(string)
	user.Email = ctx.Locals(helpers.CurrentEmail).(string)

	if err := ctx.BodyParser(&changePasswordDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	if err := a.validate.Struct(changePasswordDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	user.Password = changePasswordDto.NewPassword

	if err := a.service.ChangePassword(user, changePasswordDto.OldPassword, changePasswordDto.NewPassword); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	return ctx.JSON(
		helpers.Success(
			"Success to change password",
			nil,
			nil,
		))
}

// CreateUser implements AuthController.
func (a *authController) CurrentUser(ctx *fiber.Ctx) error {

	token := ctx.Get("Authorization")

	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				"Token is required", nil, nil,
			))
	}

	user, err := a.service.CurrentUser(strings.Replace(token, "Bearer ", "", 1))

	if err != nil {
		// check if the error contains "token is expired"
		if strings.Contains(err.Error(), "expired") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				helpers.Error(
					err.Error(), err, nil,
				))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	return ctx.JSON(
		helpers.Success(
			"Success to get current user",
			user,
			nil,
		))
}

func (a *authController) Login(ctx *fiber.Ctx) error {
	userLogin := userDto.UserLogin{}

	if err := ctx.BodyParser(&userLogin); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	err := a.validate.Struct(userLogin)

	if err != nil {

		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	token, err := a.service.Login(userLogin.Email, userLogin.Password)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			helpers.Error(
				err.Error(), err, nil,
			))
	}

	return ctx.JSON(
		helpers.Success(
			"Success to login",
			fiber.Map{
				"token":      token,
				"token_type": "Bearer",
			},
			nil,
		))
}

// Register implements AuthController.
func (a *authController) Register(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

func NewAuthController(service services.AuthService, validate *validator.Validate) AuthController {
	return &authController{
		service:  service,
		validate: validate,
	}
}
