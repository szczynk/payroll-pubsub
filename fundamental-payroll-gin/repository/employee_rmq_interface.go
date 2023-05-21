package repository

import "fundamental-payroll-gin/model"

type EmployeeRMQPubRepoI interface {
	Add(req *model.Employee) (*model.Employee, error)
}
