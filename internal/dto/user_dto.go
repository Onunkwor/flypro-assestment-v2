package dto

import "github.com/onunkwor/flypro-assestment-v2/internal/utils"

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (r *CreateUserRequest) Sanitize() {
	r.Email = utils.SanitizeString(r.Email)
	r.Name = utils.SanitizeString(r.Name)
}
