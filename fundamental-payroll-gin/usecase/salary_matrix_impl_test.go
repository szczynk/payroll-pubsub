package usecase_test

import (
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSalaryUsecase_List(t *testing.T) {
	tests := []struct {
		name       string
		repoResult []model.SalaryMatrix
		repoErr    error
		want       []model.SalaryMatrix
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			repoResult: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
			},
			want: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSalaryRepo := mocks.NewSalaryRepositoryI(t)

			uc := usecase.NewSalaryUsecase(mockSalaryRepo)

			mockSalaryRepo.On("List").Return(tt.repoResult, tt.repoErr)

			got, err := uc.List()

			t.Logf("err: %v", err)
			assert.Equal(t, tt.wantErr, err != nil, "SalaryUsecase.List() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "SalaryUsecase.List() = %v, want %v", got, tt.want)
			t.Log("\n\n")
		})
	}
}
