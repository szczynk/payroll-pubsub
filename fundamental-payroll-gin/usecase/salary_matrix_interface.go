//go:generate mockery --output=../mocks --name SalaryUsecase
package usecase

import "fundamental-payroll-gin/model"

type SalaryUsecaseI interface {
	List() ([]model.SalaryMatrix, error)
}
