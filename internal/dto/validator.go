package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// RegisterValidator 注册验证器
func RegisterValidator(validate *validator.Validate) {
	validate.RegisterValidation("year", ValidYear)
}

// ValidYear 验证年份
func ValidYear(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	currentYear := int64(time.Now().Year())
	return year >= 1900 && year <= currentYear+100 // 可根据需要设定合理年份范围
}
