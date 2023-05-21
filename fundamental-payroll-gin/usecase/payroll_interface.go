//go:generate mockery --output=../mocks --name PayrollUsecase
package usecase

import "fundamental-payroll-gin/model"

type PayrollUsecaseI interface {
	List() ([]model.Payroll, error)
	Add(req *model.PayrollRequest) (*model.Payroll, error)
	Detail(id int64) (*model.Payroll, error)
}
