//go:generate mockery --output=../mocks --name PayrollUsecase
package usecase

import "consumer-payroll/model"

type PayrollUsecaseI interface {
	Add(req *model.Payroll) (*model.Payroll, error)
}
