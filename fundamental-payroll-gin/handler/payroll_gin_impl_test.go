package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPayrollGinHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		ucResult   []model.Payroll
		ucErr      error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			ucResult: []model.Payroll{
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
			wantStatus: http.StatusOK,
		},
		{
			name:       "uc error",
			ucErr:      errors.New("uc error"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockUC := mocks.NewPayrollUsecaseI(t)

			h := handler.NewPayrollGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("List").Return(tt.ucResult, tt.ucErr)
			}

			req, errReq := http.NewRequest(http.MethodGet, "/payrolls", nil)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			c.Request = req

			h.List(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
		})
	}
}

func TestPayrollGinHandler_Add(t *testing.T) {
	type args struct {
		req *model.PayrollRequest
	}
	tests := []struct {
		name       string
		args       args
		ucResult   *model.Payroll
		ucErr      error
		wantStatus int
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
			ucResult: &model.Payroll{
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
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid employee'id",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          0,
					TotalHariMasuk:      15,
					TotalHariTidakMasuk: 5,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid total hari masuk",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          1,
					TotalHariMasuk:      -1,
					TotalHariTidakMasuk: 5,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid total hari tidak masuk",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          1,
					TotalHariMasuk:      10,
					TotalHariTidakMasuk: -1,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "uc error 1",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          1,
					TotalHariMasuk:      10,
					TotalHariTidakMasuk: 10,
				},
			},
			ucErr:      sql.ErrNoRows,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "uc error 2",
			args: args{
				req: &model.PayrollRequest{
					EmployeeID:          1,
					TotalHariMasuk:      10,
					TotalHariTidakMasuk: 10,
				},
			},
			ucErr:      errors.New("uc error 2"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockUC := mocks.NewPayrollUsecaseI(t)

			h := handler.NewPayrollGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("Add", mock.Anything).Return(tt.ucResult, tt.ucErr)
			}

			jsonBytes, _ := json.Marshal(tt.args.req)
			reqBody := bytes.NewBuffer(jsonBytes)

			req, errReq := http.NewRequest(http.MethodPost, "/payrolls", reqBody)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			c.Request = req

			h.Add(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
		})
	}
}

func TestPayrollGinHandler_Detail(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		ucResult   *model.Payroll
		ucErr      error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{id: "1"},
			ucResult: &model.Payroll{
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
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			args:       args{id: "a"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			args:       args{id: "10"},
			ucErr:      sql.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "uc error",
			args:       args{id: "10"},
			ucErr:      errors.New("uc error"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockUC := mocks.NewPayrollUsecaseI(t)

			h := handler.NewPayrollGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("Detail", mock.Anything).Return(tt.ucResult, tt.ucErr)
			}

			req, errReq := http.NewRequest(http.MethodGet, "/payrolls/"+tt.args.id, nil)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/payrolls/:id", h.Detail)
			router.ServeHTTP(res, req)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
		})
	}
}
