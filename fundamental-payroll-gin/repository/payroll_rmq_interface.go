//go:generate mockery --output=../mocks --name PayrollRMQPubRepoI
package repository

import "fundamental-payroll-gin/model"

type PayrollRMQPubRepoI interface {
	Add(req *model.Payroll) (*model.Payroll, error)
}
