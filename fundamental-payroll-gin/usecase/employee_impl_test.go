package usecase_test

import (
	"errors"
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmployeeUsecase_List(t *testing.T) {
	tests := []struct {
		name       string
		repoResult []model.Employee
		repoErr    error
		want       []model.Employee
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			repoResult: []model.Employee{
				{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
			},
			want: []model.Employee{
				{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmployeeRepo := mocks.NewEmployeeRepositoryI(t)
			mockEmployeePubRepo := mocks.NewEmployeeRMQPubRepoI(t)
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewEmployeeUsecase(
				mockEmployeeRepo,
				mockEmployeePubRepo,
				mockSalaryRepo,
			)

			mockEmployeeRepo.On("List").Return(tt.repoResult, tt.repoErr)

			got, err := uc.List()

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "EmployeeUsecase.List() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "EmployeeUsecase.List() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}

func TestEmployeeUsecase_Detail(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		repoResult *model.Employee
		repoErr    error
		want       *model.Employee
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "success",
			args:       args{id: 1},
			repoResult: &model.Employee{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
			want:       &model.Employee{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
		},
		{
			name:    "failed",
			args:    args{id: 0},
			repoErr: assert.AnError,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmployeeRepo := mocks.NewEmployeeRepositoryI(t)
			mockEmployeePubRepo := mocks.NewEmployeeRMQPubRepoI(t)
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewEmployeeUsecase(
				mockEmployeeRepo,
				mockEmployeePubRepo,
				mockSalaryRepo,
			)

			mockEmployeeRepo.On("Detail", tt.args.id).Return(tt.repoResult, tt.repoErr)

			got, err := uc.Detail(tt.args.id)

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "EmployeeUsecase.Detail() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "EmployeeUsecase.Detail() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}

func TestEmployeeUsecase_Add(t *testing.T) {
	married := true
	type args struct {
		req *model.EmployeeRequest
	}
	tests := []struct {
		name             string
		args             args
		repoSalaryResult []model.SalaryMatrix
		repoSalaryErr    error
		repoPubResult    *model.Employee
		repoPubErr       error
		want             *model.Employee
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   1,
					Married: &married,
				},
			},
			repoSalaryResult: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
			},
			repoPubResult: &model.Employee{
				Name:    "test",
				Gender:  "laki-laki",
				Grade:   1,
				Married: married,
			},
			want: &model.Employee{
				Name:    "test",
				Gender:  "laki-laki",
				Grade:   1,
				Married: married,
			},
		},
		{
			name: "invalid grade",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   4,
					Married: &married,
				},
			},
			repoSalaryResult: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
			},
			wantErr: true,
		},
		{
			name: "salary list error",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   1,
					Married: &married,
				},
			},
			repoSalaryErr: errors.New("salary list error"),
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmployeeRepo := mocks.NewEmployeeRepositoryI(t)
			mockEmployeePubRepo := mocks.NewEmployeeRMQPubRepoI(t)
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewEmployeeUsecase(
				mockEmployeeRepo,
				mockEmployeePubRepo,
				mockSalaryRepo,
			)

			mockSalaryRepo.On("List").Return(tt.repoSalaryResult, tt.repoSalaryErr)

			if tt.repoPubResult != nil || tt.repoPubErr != nil {
				mockEmployeePubRepo.On("Add", mock.Anything).Return(tt.repoPubResult, tt.repoPubErr)
			}

			got, err := uc.Add(tt.args.req)

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "EmployeeUsecase.Add() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "EmployeeUsecase.Add() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}
