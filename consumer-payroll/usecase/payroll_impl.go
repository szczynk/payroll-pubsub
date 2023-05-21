package usecase

import (
	"consumer-payroll/model"
	"consumer-payroll/repository"
)

type PayrollUsecase struct {
	payrollRepo repository.PayrollRepositoryI
}

func NewPayrollUsecase(
	payrollRepo repository.PayrollRepositoryI,
) PayrollUsecaseI {
	return &PayrollUsecase{
		payrollRepo: payrollRepo,
	}
}

func (uc *PayrollUsecase) Add(payroll *model.Payroll) (*model.Payroll, error) {
	return uc.payrollRepo.Add(payroll)
}
