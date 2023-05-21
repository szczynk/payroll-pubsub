package repository

import (
	"fundamental-payroll-gin/helper/timeout"
	"fundamental-payroll-gin/model"

	"gorm.io/gorm"
)

type PayrollGormRepository struct {
	db *gorm.DB
}

func NewPayrollGormRepository(db *gorm.DB) PayrollRepositoryI {
	r := new(PayrollGormRepository)
	r.db = db
	return r
}

func (repo *PayrollGormRepository) List() ([]model.Payroll, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var payroll model.Payroll
	var payrolls []model.Payroll

	sqlQuery := `
	SELECT payrolls.id, payrolls.basic_salary, payrolls.pay_cut, payrolls.additional_salary, payrolls.employee_id,
				employees.id, employees.name, employees.gender, employees.grade, employees.married
	FROM payrolls
	INNER JOIN employees ON payrolls.employee_id = employees.id
	ORDER BY payrolls.id ASC
	`
	rows, sqlErr := repo.db.WithContext(ctx).Raw(sqlQuery).Rows()
	if sqlErr != nil {
		return payrolls, sqlErr
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(
			&payroll.ID, &payroll.BasicSalary, &payroll.PayCut, &payroll.AdditionalSalary, &payroll.EmployeeID,
			&payroll.Employee.ID, &payroll.Employee.Name, &payroll.Employee.Gender, &payroll.Employee.Grade,
			&payroll.Employee.Married,
		)
		if scanErr != nil {
			return payrolls, scanErr
		}

		payrolls = append(payrolls, payroll)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return payrolls, rowErr
	}
	payroll = model.Payroll{}

	return payrolls, nil
}

func (repo *PayrollGormRepository) Detail(id int64) (*model.Payroll, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	SELECT payrolls.id, payrolls.basic_salary, payrolls.pay_cut, payrolls.additional_salary, payrolls.employee_id, 
				 employees.id, employees.name, employees.gender, employees.grade, employees.married 
	FROM payrolls 
	INNER JOIN employees ON payrolls.employee_id = employees.id 
	WHERE payrolls.id = $1 LIMIT 1
	`
	row := repo.db.WithContext(ctx).Raw(sqlQuery, id).Row()

	payroll := new(model.Payroll)
	scanErr := row.Scan(
		&payroll.ID, &payroll.BasicSalary, &payroll.PayCut, &payroll.AdditionalSalary, &payroll.EmployeeID,
		&payroll.Employee.ID, &payroll.Employee.Name, &payroll.Employee.Gender, &payroll.Employee.Grade,
		&payroll.Employee.Married,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return payroll, nil
}
