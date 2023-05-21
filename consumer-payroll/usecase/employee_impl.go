package usecase

import (
	"consumer-payroll/model"
	"consumer-payroll/repository"
)

type EmployeeUsecase struct {
	employeeRepo repository.EmployeeRepositoryI
}

func NewEmployeeUsecase(
	employeeRepo repository.EmployeeRepositoryI,
) EmployeeUsecaseI {
	return &EmployeeUsecase{
		employeeRepo: employeeRepo,
	}
}

func (uc *EmployeeUsecase) Add(employee *model.Employee) (*model.Employee, error) {
	return uc.employeeRepo.Add(employee)
}
