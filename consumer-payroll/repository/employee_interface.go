//go:generate mockery --output=../mocks --name EmployeeRepository
package repository

import "consumer-payroll/model"

type EmployeeRepositoryI interface {
	Add(req *model.Employee) (*model.Employee, error)
	Detail(id int64) (*model.Employee, error)
}
