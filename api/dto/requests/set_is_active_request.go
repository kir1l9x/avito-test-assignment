package requests

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required,user_id"`
	IsActive bool   `json:"is_active"`
}
