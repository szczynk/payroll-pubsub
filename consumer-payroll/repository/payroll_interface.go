//go:generate mockery --output=../mocks --name PayrollRepository
package repository

import "consumer-payroll/model"

type PayrollRepositoryI interface {
	Add(payroll *model.Payroll) (*model.Payroll, error)
}
