package handlers

import (
	"siska-rgb-golang-test/internal/helpers/jwt"
	uuid_helpers "siska-rgb-golang-test/internal/helpers/uuid"
	helpers "siska-rgb-golang-test/internal/helpers/validator"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/users/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HandlerUser struct {
	userService services.Service
}

func NewHandlerUser(userService services.Service) *HandlerUser {
	return &HandlerUser{userService: userService}
}

func (h *HandlerUser) Register(c *fiber.Ctx) error {

	logger := c.Locals("logger").(*logrus.Logger)

	var request models.UserRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	isEmailExists, _ := h.userService.IsEmailExists(c.Context(), request.Email)
	if isEmailExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Email",
			"message": "Email already exists"},
		})
	}

	userResponse, err := h.userService.RegisterUser(c.Context(), request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method":  c.Method(),
			"route":   c.Path(),
			"error":   err.Error(),
			"payload": request,
		}).Error("Failed to parse register user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(userResponse)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to generate token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	result := models.UserResponseWithToken{
		UserResponse: userResponse,
		Token:        token,
	}

	logger.WithFields(logrus.Fields{
		"method": c.Method(),
		"route":  c.Path(),
		"error":  nil,
	}).Info("Success login user")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (h *HandlerUser) Login(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	var request models.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	userResponse, err := h.userService.LoginUser(c.Context(), request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to parse login user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Email or password is incorrect"})
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(userResponse)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to generate token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	result := models.UserResponseWithToken{
		UserResponse: userResponse,
		Token:        token,
	}

	logger.WithFields(logrus.Fields{
		"method": c.Method(),
		"route":  c.Path(),
		"error":  nil,
	}).Info("Success login user")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

func (h *HandlerUser) Profile(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	id := uuid_helpers.ParseUUID(c.Locals("user").(*models.UserResponse).ID)

	isUserIDExists, err := h.userService.IsUserIDExists(c.Context(), id)
	if !isUserIDExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	userResponse, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get user by id")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": userResponse})
}

func (h *HandlerUser) Update(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	id := uuid_helpers.ParseUUID(c.Locals("user").(*models.UserResponse).ID)

	isUserIDExists, err := h.userService.IsUserIDExists(c.Context(), id)
	if !isUserIDExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	var request models.UserRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	if request.Email != c.Locals("user").(*models.UserResponse).Email {
		isEmailExists, _ := h.userService.IsEmailExists(c.Context(), request.Email)
		if isEmailExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
				"field":   "Email",
				"message": "Email already exists"},
			})
		}
	}

	userResponse, err := h.userService.UpdateUser(c.Context(), id, request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to update user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	c.Locals("user", userResponse)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": userResponse})
}
