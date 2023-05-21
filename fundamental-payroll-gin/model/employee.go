package model

type Employee struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Grade   int8   `json:"grade"`
	Married bool   `json:"married"`
}

type EmployeeRequest struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Gender  string `form:"gender" json:"gender" binding:"required"`
	Grade   int8   `form:"grade" json:"grade" binding:"required"`
	Married *bool  `form:"married" json:"married" binding:"required"`
}
