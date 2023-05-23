//go:generate mockery --output=../mocks --name SalaryUsecaseI
package usecase

import "fundamental-payroll-gin/model"

type SalaryUsecaseI interface {
	List() ([]model.SalaryMatrix, error)
}
