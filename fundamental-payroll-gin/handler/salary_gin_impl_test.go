package handler_test

import (
	"errors"
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSalaryGinHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		ucResult   []model.SalaryMatrix
		ucErr      error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			ucResult: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
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
			mockUC := mocks.NewSalaryUsecaseI(t)

			h := handler.NewSalaryGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("List").Return(tt.ucResult, tt.ucErr)
			}

			req, errReq := http.NewRequest(http.MethodGet, "/salaries", nil)
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
