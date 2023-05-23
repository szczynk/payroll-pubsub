package usecase_test

import (
	"database/sql"
	"errors"
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPayrollUsecase_List(t *testing.T) {
	tests := []struct {
		name       string
		repoResult []model.Payroll
		repoErr    error
		want       []model.Payroll
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			repoResult: []model.Payroll{
				{
					ID:               1,
					BasicSalary:      5000000,
					PayCut:           400000,
					AdditionalSalary: 3400000,
					EmployeeID:       1,
					Employee: model.Employee{
						ID:      1,
						Name:    "test1",
						Gender:  "laki-laki",
						Grade:   1,
						Married: true,
					},
				},
			},
			want: []model.Payroll{
				{
					ID:               1,
					BasicSalary:      5000000,
					PayCut:           400000,
					AdditionalSalary: 3400000,
					EmployeeID:       1,
					Employee: model.Employee{
						ID:      1,
						Name:    "test1",
						Gender:  "laki-laki",
						Grade:   1,
						Married: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPayrollRepo := mocks.NewPayrollRepositoryI(t)
			mockPayrollPubRepo := mocks.NewPayrollRMQPubRepoI(t)
			mockEmployeeRepo := mocks.NewEmployeeRepositoryI(t)
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewPayrollUsecase(
				mockPayrollRepo,
				mockPayrollPubRepo,
				mockEmployeeRepo,
				mockSalaryRepo,
			)

			mockPayrollRepo.On("List").Return(tt.repoResult, tt.repoErr)

			got, err := uc.List()

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "PayrollUsecase.List() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "PayrollUsecase.List() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}

func TestPayrollUsecase_Detail(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		repoResult *model.Payroll
		repoErr    error
		want       *model.Payroll
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{id: 1},
			repoResult: &model.Payroll{
				ID:               1,
				BasicSalary:      5000000,
				PayCut:           400000,
				AdditionalSalary: 3400000,
				EmployeeID:       1,
				Employee: model.Employee{
					ID:      1,
					Name:    "test1",
					Gender:  "laki-laki",
					Grade:   1,
					Married: true,
				},
			},
			want: &model.Payroll{
				ID:               1,
				BasicSalary:      5000000,
				PayCut:           400000,
				AdditionalSalary: 3400000,
				EmployeeID:       1,
				Employee: model.Employee{
					ID:      1,
					Name:    "test1",
					Gender:  "laki-laki",
					Grade:   1,
					Married: true,
				},
			},
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
			mockPayrollRepo := mocks.NewPayrollRepositoryI(t)
			mockPayrollPubRepo := mocks.NewPayrollRMQPubRepoI(t)
			mockEmployeeRepo := mocks.NewEmployeeRepositoryI(t)
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewPayrollUsecase(
				mockPayrollRepo,
				mockPayrollPubRepo,
				mockEmployeeRepo,
				mockSalaryRepo,
			)

			mockPayrollRepo.On("Detail", tt.args.id).Return(tt.repoResult, tt.repoErr)

			got, err := uc.Detail(tt.args.id)

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "PayrollUsecase.Detail() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "PayrollUsecase.Detail() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}

func TestPayrollUsecase_Add(t *testing.T) {
	type args struct {
		req *model.PayrollRequest
	}
	tests := []struct {
		name               string
		args               args
		repoEmployeeResult *model.Employee
		repoEmployeeErr    error
		repoSalaryResult   []model.SalaryMatrix
		repoSalaryErr      error
		repoPubResult      *model.Payroll
		repoPubErr         error
		want               *model.Payroll
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          1,
					TotalHariMasuk:      15,
					TotalHariTidakMasuk: 5,
				},
			},
			repoEmployeeResult: &model.Employee{
				ID:      1,
				Name:    "test",
				Gender:  "laki-laki",
				Grade:   1,
				Married: true,
			},
			repoSalaryResult: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
			},
			repoPubResult: &model.Payroll{
				BasicSalary:      5000000,
				PayCut:           500000,
				AdditionalSalary: 3250000,
				Employee: model.Employee{
					ID:      1,
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   1,
					Married: true,
				},
				EmployeeID: 1,
			},
			want: &model.Payroll{
				BasicSalary:      5000000,
				PayCut:           500000,
				AdditionalSalary: 3250000,
				Employee: model.Employee{
					ID:      1,
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   1,
					Married: true,
				},
				EmployeeID: 1,
			},
		},
		{
			name: "invalid employee's id",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          10,
					TotalHariMasuk:      15,
					TotalHariTidakMasuk: 5,
				},
			},
			repoEmployeeErr: sql.ErrNoRows,
			wantErr:         true,
		},
		{
			name: "salary list error",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          1,
					TotalHariMasuk:      15,
					TotalHariTidakMasuk: 5,
				},
			},
			repoEmployeeResult: &model.Employee{
				ID:      1,
				Name:    "test",
				Gender:  "laki-laki",
				Grade:   1,
				Married: true,
			},
			repoSalaryErr: errors.New("salary list error"),
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPayrollRepo := mocks.NewPayrollRepositoryI(t)
			mockPayrollPubRepo := mocks.NewPayrollRMQPubRepoI(t)
			mockEmployeeRepo := mocks.NewEmployeeRepositoryI(t)
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewPayrollUsecase(
				mockPayrollRepo,
				mockPayrollPubRepo,
				mockEmployeeRepo,
				mockSalaryRepo,
			)

			if tt.repoEmployeeResult != nil || tt.repoEmployeeErr != nil {
				mockEmployeeRepo.On("Detail", tt.args.req.EmployeeID).Return(tt.repoEmployeeResult, tt.repoEmployeeErr)
			}

			if tt.repoSalaryResult != nil || tt.repoSalaryErr != nil {
				mockSalaryRepo.On("List").Return(tt.repoSalaryResult, tt.repoSalaryErr)
			}

			if tt.repoPubResult != nil || tt.repoPubErr != nil {
				mockPayrollPubRepo.On("Add", mock.Anything).Return(tt.repoPubResult, tt.repoPubErr)
			}

			got, err := uc.Add(tt.args.req)

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "PayrollUsecase.Add() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "PayrollUsecase.Add() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}
