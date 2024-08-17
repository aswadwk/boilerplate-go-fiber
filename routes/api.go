package routes

import (
	"fmt"

	"aswadwk/chatai/config"
	"aswadwk/chatai/controllers"
	"aswadwk/chatai/helpers"
	"aswadwk/chatai/middleware"
	"aswadwk/chatai/repositories"
	"aswadwk/chatai/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	db, _    = config.DBConnect()
	sqlDB, _ = db.DB()
	validate = validator.New(validator.WithRequiredStructEnabled())
	_        = validate.RegisterValidation("date_birth", helpers.ValidateDateBirthFormat)
	_        = validate.RegisterValidation("array_uuid", helpers.ValidateArrayString)

	// User
	userRepo    repositories.UserRepository = repositories.NewUserRepository(db)
	userService services.UserService        = services.NewUserService(userRepo)
	userCtrl    controllers.UserController  = controllers.NewUserController(userService, validate)

	// JWT
	// jwtRepo    repositories.JwtRepository = repositories.NewJwtRepository(db)
	jwtService services.JwtService = services.NewJwtService(userRepo)

	// Middleware
	userMiddleware middleware.AuthMiddleware = middleware.NewAuthService(jwtService)

	// Auth
	authService services.AuthService       = services.NewAuthService(userService, jwtService)
	authCtrl    controllers.AuthController = controllers.NewAuthController(authService, validate)
)

func Setup(app *fiber.App) {
	err := db.Error
	if err != nil {
		fmt.Println("Error Database !")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		stash := config.CheckIdleConnections(sqlDB)
		// jobs, err := jobRepo.GetJobsByQueue(models.JobExtractFile)

		if err != nil {
			return nil
		}

		return c.JSON(fiber.Map{
			"Database": stash,
			"Session":  fiber.Map{
				// "jobs": jobs,
			},
		})
	})

	// access assert
	app.Get("/assets/*", func(c *fiber.Ctx) error {
		return c.SendFile(c.Params("*", "public"))
	})

	api := app.Group("/api/v1")
	{
		api.Get("/tes", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World 👋!")
		})
	}

	// Auth
	auth := api.Group("/auth")
	{
		auth.Post("/login", authCtrl.Login)
		auth.Use(userMiddleware.Protected()).Get("/current-user", authCtrl.CurrentUser)
		auth.Use(userMiddleware.Protected()).Post("/change-password", authCtrl.ChangePassword)
		// auth.Use(userMiddleware.Protected()).Post("/users", userCtrl.SaveUser)
		// auth.Use(userMiddleware.Protected()).Get("/users", userCtrl.GetUsers)
		// auth.Use(userMiddleware.Protected()).Delete("/users/:id", userCtrl.DeleteUser)
		// auth.Use(userMiddleware.Protected()).Delete("/user/:by/:value", userCtrl.FindUserBy)
	}

}
