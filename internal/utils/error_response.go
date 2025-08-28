package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	out := map[string]string{}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range errs {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				out[field] = fmt.Sprintf("%s is required", field)
			case "email":
				out[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				out[field] = fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
			case "max":
				out[field] = fmt.Sprintf("%s must be at most %s characters long", field, fe.Param())
			case "len":
				out[field] = fmt.Sprintf("%s must be exactly %s characters long", field, fe.Param())
			case "oneof":
				out[field] = fmt.Sprintf("%s must be one of the following: %s", field, fe.Param())
			default:
				out[field] = fmt.Sprintf("%s is not valid (%s)", field, fe.Tag())
			}
		}
	} else {
		out["error"] = err.Error()
	}
	return out
}

func ValidationErrorResponse(c *gin.Context, err map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{"validation_error": gin.H{"details": err}})
}

func InternalServerErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "internal server error",
		"details": err.Error(),
	})
}

func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "bad_request",
		"message": message,
	})
}

func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"error":   "not_found",
		"message": message,
	})
}

func DuplicateEntryResponse(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, gin.H{
		"error":   "duplicate_entry",
		"message": message,
	})
}
