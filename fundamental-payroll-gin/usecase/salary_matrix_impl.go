package usecase

import (
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/repository"
)

type SalaryUsecase struct {
	salaryRepo repository.SalaryRepositoryI
}

func NewSalaryUsecase(salaryRepo repository.SalaryRepositoryI) SalaryUsecaseI {
	return &SalaryUsecase{
		salaryRepo: salaryRepo,
	}
}

func (uc *SalaryUsecase) List() ([]model.SalaryMatrix, error) {
	return uc.salaryRepo.List()
}
