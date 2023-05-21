package repository

import (
	"fundamental-payroll-gin/helper/timeout"
	"fundamental-payroll-gin/model"

	"gorm.io/gorm"
)

type EmployeeGormRepository struct {
	db *gorm.DB
}

func NewEmployeeGormRepository(db *gorm.DB) EmployeeRepositoryI {
	r := new(EmployeeGormRepository)
	r.db = db
	return r
}

func (repo *EmployeeGormRepository) List() ([]model.Employee, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var employee model.Employee
	var employees []model.Employee

	sqlQuery := "SELECT id, name, gender, grade, married FROM employees ORDER BY id ASC"
	rows, sqlErr := repo.db.WithContext(ctx).Raw(sqlQuery).Rows()
	if sqlErr != nil {
		return employees, sqlErr
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(&employee.ID, &employee.Name, &employee.Gender, &employee.Grade, &employee.Married)
		if scanErr != nil {
			return employees, scanErr
		}

		employees = append(employees, employee)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return employees, rowErr
	}
	employee = model.Employee{}

	return employees, nil
}

func (repo *EmployeeGormRepository) Detail(id int64) (*model.Employee, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := "SELECT id, name, gender, grade, married FROM employees WHERE id = $1 LIMIT 1"
	row := repo.db.WithContext(ctx).Raw(sqlQuery, id).Row()

	employee := new(model.Employee)
	scanErr := row.Scan(&employee.ID, &employee.Name, &employee.Gender, &employee.Grade, &employee.Married)
	if scanErr != nil {
		return nil, scanErr
	}

	return employee, nil
}
