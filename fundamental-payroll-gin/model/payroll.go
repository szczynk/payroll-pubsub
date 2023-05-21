package model

type Payroll struct {
	ID               int64    `json:"id"`
	BasicSalary      int64    `json:"basic_salary"`
	PayCut           int64    `json:"pay_cut"`
	AdditionalSalary int64    `json:"additional_salary"`
	Employee         Employee `json:"employee"`
	EmployeeID       int64    `json:"employee_id"`
}

type PayrollRequest struct {
	EmployeeID          int64 `form:"employee_id" json:"employee_id" binding:"required"`
	TotalHariMasuk      int64 `form:"total_hari_masuk" json:"total_hari_masuk" binding:"required"`
	TotalHariTidakMasuk int64 `form:"total_hari_tidak_masuk" json:"total_hari_tidak_masuk" binding:"required"`
}
