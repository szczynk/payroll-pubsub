//go:generate mockery --output=../mocks --name PayrollRepositoryI
package repository

import "fundamental-payroll-gin/model"

type PayrollRepositoryI interface {
	List() ([]model.Payroll, error)
	Detail(id int64) (*model.Payroll, error)
}
