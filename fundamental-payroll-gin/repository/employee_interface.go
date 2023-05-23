//go:generate mockery --output=../mocks --name EmployeeRepositoryI
package repository

import "fundamental-payroll-gin/model"

type EmployeeRepositoryI interface {
	List() ([]model.Employee, error)
	Detail(id int64) (*model.Employee, error)
}
