//go:generate mockery --output=../mocks --name SalaryRepositoryI
package repository

import "fundamental-payroll-gin/model"

type SalaryRepositoryI interface {
	List() ([]model.SalaryMatrix, error)
}
