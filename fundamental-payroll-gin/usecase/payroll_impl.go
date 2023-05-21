package usecase

import (
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/repository"
	"strings"
)

type PayrollUsecase struct {
	payrollRepo  repository.PayrollRepositoryI
	payrollPub   repository.PayrollRMQPubRepoI
	employeeRepo repository.EmployeeRepositoryI
	salaryRepo   repository.SalaryRepositoryI
}

func NewPayrollUsecase(
	payrollRepo repository.PayrollRepositoryI,
	payrollPub repository.PayrollRMQPubRepoI,
	employeeRepo repository.EmployeeRepositoryI,
	salaryRepo repository.SalaryRepositoryI,
) PayrollUsecaseI {
	return &PayrollUsecase{
		payrollRepo:  payrollRepo,
		payrollPub:   payrollPub,
		employeeRepo: employeeRepo,
		salaryRepo:   salaryRepo,
	}
}

func (uc *PayrollUsecase) List() ([]model.Payroll, error) {
	return uc.payrollRepo.List()
}

// note: tldr
// to be fair, is hard applying usecase inside consumer-payroll
// if one of employee's or salaries' repo returning error
// by biz logic, i should send nil or null or empty string
// from consumer-payroll's publisher across rabbitmq,
// make fundamental-payroll-gin's consumer detect it,
// and finally returning error response
// but how supposed to that? yes ChatGPT dumb if too much context
// this is the way i can use right now
// i'm suck at 2 service pubsub each other.
func (uc *PayrollUsecase) Add(req *model.PayrollRequest) (*model.Payroll, error) {
	employee, err := uc.employeeRepo.Detail(req.EmployeeID)
	if err != nil {
		return nil, err
	}

	var (
		basicSalary      int64
		payCut           int64
		additionalSalary int64
	)

	salaries, err := uc.salaryRepo.List()
	if err != nil {
		return nil, err
	}

	for _, salary := range salaries {
		if salary.Grade == employee.Grade {
			basicSalary = salary.BasicSalary
			payCut = salary.PayCut * req.TotalHariTidakMasuk

			additionalSalary = salary.Allowance * req.TotalHariMasuk
			if strings.Contains(strings.ToLower(employee.Gender), "laki-laki") && employee.Married {
				additionalSalary += salary.HoF
			}
		}
	}

	payroll := &model.Payroll{
		BasicSalary:      basicSalary,
		PayCut:           payCut,
		AdditionalSalary: additionalSalary,
		Employee:         *employee,
		EmployeeID:       req.EmployeeID,
	}

	return uc.payrollPub.Add(payroll)
}

func (uc *PayrollUsecase) Detail(id int64) (*model.Payroll, error) {
	return uc.payrollRepo.Detail(id)
}
