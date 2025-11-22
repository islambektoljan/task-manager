package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var validate *validator.Validate

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v

		// Регистрируем кастомные валидации
		validate.RegisterValidation("notblank", validators.NotBlank)

		// Можно добавить свои кастомные валидации
		validate.RegisterValidation("strongpassword", validateStrongPassword)
	}
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Проверка сложности пароля
	if len(password) < 8 {
		return false
	}

	// Дополнительные проверки можно добавить
	return true
}

func GetValidator() *validator.Validate {
	return validate
}
