package repository

import (
	"consumer-payroll/helper/timeout"
	"consumer-payroll/model"

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

func (repo *PayrollGormRepository) Add(payroll *model.Payroll) (*model.Payroll, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	WITH inserted_payroll AS (
		INSERT INTO payrolls (basic_salary, pay_cut, additional_salary, employee_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, basic_salary, pay_cut, additional_salary, employee_id
	)
	SELECT inserted_payroll.id, inserted_payroll.basic_salary, inserted_payroll.pay_cut, 
				 inserted_payroll.additional_salary, inserted_payroll.employee_id, 
				 employees.id, employees.name, employees.gender, employees.grade, employees.married
	FROM inserted_payroll
	INNER JOIN employees ON inserted_payroll.employee_id = employees.id
	`
	row := repo.db.WithContext(ctx).Raw(sqlQuery, payroll.BasicSalary, payroll.PayCut,
		payroll.AdditionalSalary, payroll.EmployeeID).Row()

	newPayroll := new(model.Payroll)
	scanErr := row.Scan(
		&newPayroll.ID, &newPayroll.BasicSalary, &newPayroll.PayCut, &newPayroll.AdditionalSalary, &newPayroll.EmployeeID,
		&newPayroll.Employee.ID, &newPayroll.Employee.Name, &newPayroll.Employee.Gender, &newPayroll.Employee.Grade,
		&newPayroll.Employee.Married,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return newPayroll, nil
}
