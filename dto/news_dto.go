package dto

type NewsInput struct {
	Title       string `json:"title" binding:"required,max=50"`
	Description string `json:"description" binding:"required,max=100"`
}
