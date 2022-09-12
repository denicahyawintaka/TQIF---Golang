package dto

// DiaryUpdateDTO is a model that client use when updating a diary
type DiaryUpdateDTO struct {
	ID       uint64 `json:"id" form:"id" binding:"required"`
	Body     string `json:"body" form:"body" binding:"required"`
	Mood     string `json:"mood" form:"mood" binding:"required"`
	Category string `json:"category" form:"category" binding:"required"`
	Date     string `json:"date" form:"date" binding:"required"`
	UserID   uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}

// DiaryCreateDTO is is a model that clinet use when create a new diary
type DiaryCreateDTO struct {
	Body     string `json:"body" form:"body" binding:"required"`
	Mood     string `json:"mood" form:"mood" binding:"required"`
	Category string `json:"category" form:"category" binding:"required"`
	Date     string `json:"date" form:"date" binding:"required"`
	UserID   uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}
