package todo

import "time"

type Task struct {
	ID          int        `json:"id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date" swaggertype:"string" example:"2024-06-07T15:00:00Z"`
	Completed   bool       `json:"completed"`
}
type Pages struct {
	CountPage int     `json:"count_page"`
	CurPage   int     `json:"cur_page"`
	Tasks     []*Task `json:"tasks"`
}
type ErrorResponse struct {
	Message string `swaggertype:"string" example:"Error"`
}
