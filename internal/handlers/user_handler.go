package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/dto"
	"github.com/onunkwor/flypro-assestment-v2/internal/models"
	"github.com/onunkwor/flypro-assestment-v2/internal/services"
	"github.com/onunkwor/flypro-assestment-v2/internal/utils"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
	request.Sanitize()
	user := models.User{
		Email: request.Email,
		Name:  request.Name,
	}
	if err := h.service.CreateUser(c.Request.Context(), &user); err != nil {
		if err == services.ErrEmailAlreadyExists {
			utils.DuplicateEntryResponse(c, "email already exists")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid user ID")
		return
	}
	response, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.NotFoundResponse(c, "user not found")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": response})
}
