package repository

import (
	"consumer-payroll/helper/timeout"
	"consumer-payroll/model"

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

func (repo *EmployeeGormRepository) Add(employee *model.Employee) (*model.Employee, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	INSERT INTO employees (name, gender, grade, married) VALUES ($1, $2, $3, $4) 
	RETURNING id, name, gender, grade, married
	`
	row := repo.db.WithContext(ctx).Raw(sqlQuery, employee.Name, employee.Gender, employee.Grade, employee.Married).Row()

	newEmployee := new(model.Employee)
	scanErr := row.Scan(&newEmployee.ID, &newEmployee.Name, &newEmployee.Gender, &newEmployee.Grade, &newEmployee.Married)
	if scanErr != nil {
		return nil, scanErr
	}

	return newEmployee, nil
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
