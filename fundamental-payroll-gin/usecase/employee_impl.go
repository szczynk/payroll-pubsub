package usecase

import (
	"errors"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/repository"
)

type EmployeeUsecase struct {
	employeeRepo repository.EmployeeRepositoryI
	employeePub  repository.EmployeeRMQPubRepoI
	salaryRepo   repository.SalaryRepositoryI
}

func NewEmployeeUsecase(
	employeeRepo repository.EmployeeRepositoryI,
	employeePub repository.EmployeeRMQPubRepoI,
	salaryRepo repository.SalaryRepositoryI,
) EmployeeUsecaseI {
	return &EmployeeUsecase{
		employeeRepo: employeeRepo,
		employeePub:  employeePub,
		salaryRepo:   salaryRepo,
	}
}

func (uc *EmployeeUsecase) List() ([]model.Employee, error) {
	return uc.employeeRepo.List()
}

func (uc *EmployeeUsecase) Add(req *model.EmployeeRequest) (*model.Employee, error) {
	salaries, err := uc.salaryRepo.List()
	if err != nil {
		return nil, err
	}

	var isValidGrade bool
	for _, salary := range salaries {
		if req.Grade == salary.Grade {
			isValidGrade = true
			break
		}
	}
	if !isValidGrade {
		return nil, errors.New("invalid employee's salary grade")
	}

	employee := &model.Employee{
		Name:    req.Name,
		Gender:  req.Gender,
		Grade:   req.Grade,
		Married: *req.Married,
	}

	return uc.employeePub.Add(employee)
}

func (uc *EmployeeUsecase) Detail(id int64) (*model.Employee, error) {
	return uc.employeeRepo.Detail(id)
}
