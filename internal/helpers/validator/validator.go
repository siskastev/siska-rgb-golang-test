package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type validationErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email format."
	case "min":
		return "This field must be at least " + fe.Param() + " characters long."
	case "max":
		return "This field must be at most " + fe.Param() + " characters long."
	default:
		return fe.Error()
	}
}

func ValidationErrorResponse(c *fiber.Ctx, errs validator.ValidationErrors) error {
	var validationErrors []validationErrors
	for _, err := range errs {
		validationMessage := getErrorMessage(err)

		validationErrors = append(validationErrors, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{
			Field:   err.Field(),
			Message: validationMessage,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"errors": validationErrors,
	})
}
