package repository

import (
	"fundamental-payroll-gin/helper/timeout"
	"fundamental-payroll-gin/model"

	"gorm.io/gorm"
)

type SalaryGormRepository struct {
	db *gorm.DB
}

func NewSalaryGormRepository(db *gorm.DB) SalaryRepositoryI {
	r := new(SalaryGormRepository)
	r.db = db
	return r
}

func (repo *SalaryGormRepository) List() ([]model.SalaryMatrix, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var salary model.SalaryMatrix
	var salaries []model.SalaryMatrix

	sqlQuery := "SELECT id, grade, basic_salary, pay_cut, allowance, head_of_family FROM salaries ORDER BY id ASC"
	rows, sqlErr := repo.db.WithContext(ctx).Raw(sqlQuery).Rows()
	if sqlErr != nil {
		return salaries, sqlErr
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(&salary.ID, &salary.Grade, &salary.BasicSalary, &salary.PayCut, &salary.Allowance, &salary.HoF)
		if scanErr != nil {
			return salaries, scanErr
		}

		salaries = append(salaries, salary)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return salaries, rowErr
	}
	salary = model.SalaryMatrix{}

	return salaries, nil
}
