//go:generate mockery --output=../mocks --name EmployeeUsecase
package usecase

import "fundamental-payroll-gin/model"

type EmployeeUsecaseI interface {
	List() ([]model.Employee, error)
	Add(req *model.EmployeeRequest) (*model.Employee, error)
	Detail(id int64) (*model.Employee, error)
}
