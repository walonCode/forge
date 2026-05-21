package tasks

import "time"


type Task struct {
	ID string `json:"id"`
	Title string `json:"task"`
	UserId string `json:"userId"`
	Description string `json:"description"`
	IsCompleted bool `json:"is_completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//create task request
type CreateTaskRequest struct {
	Title string `json:"title" validate:"min=2,required"`
	Description string `json:"description" validate:"min=2,required"`
}

//create task response 
type  TaskResponse struct {
	Message string `json:"message"`
	Data any `json:"data"` 
}




//database
type CreateTaskParam struct {
	ID string
	Title string 
	Description string
	UserId string
	IsCompleted bool
}