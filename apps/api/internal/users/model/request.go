package model

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"  binding:"required,min=1,max=128"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"omitempty,min=1,max=128"`
	Plan string `json:"plan" binding:"omitempty,oneof=free pro enterprise"`
}
