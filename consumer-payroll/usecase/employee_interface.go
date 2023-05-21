//go:generate mockery --output=../mocks --name EmployeeUsecase
package usecase

import "consumer-payroll/model"

type EmployeeUsecaseI interface {
	Add(req *model.Employee) (*model.Employee, error)
}
