package model

type SalaryMatrix struct {
	ID          int64 `json:"id"`
	Grade       int8  `json:"grade"`
	BasicSalary int64 `json:"basic_salary"`
	PayCut      int64 `json:"pay_cut"`
	Allowance   int64 `json:"allowance"`
	HoF         int64 `json:"head_of_family"`
}

type SalaryMatrixRequest struct {
	Grade       int8  `form:"grade" json:"grade" binding:"required"`
	BasicSalary int64 `form:"basic_salary" json:"basic_salary" binding:"required"`
	PayCut      int64 `form:"pay_cut" json:"pay_cut" binding:"required"`
	Allowance   int64 `form:"allowance" json:"allowance" binding:"required"`
	HoF         int64 `form:"head_of_family" json:"head_of_family" binding:"required"`
}
