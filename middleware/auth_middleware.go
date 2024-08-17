package middleware

import (
	"strings"

	"aswadwk/chatai/helpers"
	"aswadwk/chatai/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware interface {
	IsSuperAdmin() fiber.Handler
	Protected() fiber.Handler
}

type authMiddleware struct {
	jwtService services.JwtService
}

// IsSuperAdmin implements AuthMiddleware.
func (a *authMiddleware) IsSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.Error(
					"Token is required", nil, nil,
				),
			)
		}

		tokenParse, err := a.jwtService.ValidateToken(strings.Replace(token, "Bearer ", "", 1))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.Error(
					"Token invalid or expired", err, nil,
				),
			)
		}

		var email, role, userId string

		for key, value := range tokenParse.Claims.(jwt.MapClaims) {
			if key == "sub" {
				email = value.(string)
			}

			if key == "jti" {
				userId = value.(string)
			}

			if key == "role" {
				role = value.(string)
			}

		}

		c.Locals(helpers.CurrentEmail, email)
		c.Locals(helpers.CurrentRole, role)
		c.Locals(helpers.CurrentUserID, userId)

		if role != helpers.RoleSuperAdmin {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.Error(
					"Not authorized", nil, nil,
				),
			)
		}

		return c.Next()
	}

}

// Protected implements AuthMiddleware.
func (a *authMiddleware) Protected() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.Error(
					"Token is required", nil, nil,
				),
			)
		}

		tokenParse, err := a.jwtService.ValidateToken(strings.Replace(token, "Bearer ", "", 1))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.Error(
					"Token invalid or expired", err, nil,
				),
			)
		}

		var email, role, userId string

		for key, value := range tokenParse.Claims.(jwt.MapClaims) {
			if key == "sub" {
				email = value.(string)
			}

			if key == "jti" {
				userId = value.(string)
			}

			if key == "role" {
				role = value.(string)
			}
		}

		c.Locals(helpers.CurrentEmail, email)
		c.Locals(helpers.CurrentRole, role)
		c.Locals(helpers.CurrentUserID, userId)

		return c.Next()
	}
}

func NewAuthService(jwtService services.JwtService) AuthMiddleware {
	return &authMiddleware{
		jwtService: jwtService,
	}
}
