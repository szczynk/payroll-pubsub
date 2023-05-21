//go:generate mockery --output=../mocks --name SalaryRepository
package repository

import "fundamental-payroll-gin/model"

type SalaryRepositoryI interface {
	List() ([]model.SalaryMatrix, error)
}
