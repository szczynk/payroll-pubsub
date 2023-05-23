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

func TestEmployeeGinHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		ucResult   []model.Employee
		ucErr      error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			ucResult: []model.Employee{
				{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
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
			mockUC := mocks.NewEmployeeUsecaseI(t)

			h := handler.NewEmployeeGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("List").Return(tt.ucResult, tt.ucErr)
			}

			req, errReq := http.NewRequest(http.MethodGet, "/employees", nil)
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

func TestEmployeeGinHandler_Add(t *testing.T) {
	married := false
	type args struct {
		req *model.EmployeeRequest
	}
	tests := []struct {
		name       string
		args       args
		ucResult   *model.Employee
		ucErr      error
		wantStatus int
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
			ucResult: &model.Employee{
				ID:      1,
				Name:    "test",
				Gender:  "laki-laki",
				Grade:   1,
				Married: married,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid field",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "",
					Gender:  "laki-laki",
					Grade:   1,
					Married: &married,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid gender",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "test",
					Gender:  "transgender",
					Grade:   1,
					Married: &married,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid grade",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   -1,
					Married: &married,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "uc error",
			args: args{
				req: &model.EmployeeRequest{
					Name:    "test",
					Gender:  "laki-laki",
					Grade:   1,
					Married: &married,
				},
			},
			ucErr:      errors.New("uc error"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockUC := mocks.NewEmployeeUsecaseI(t)

			h := handler.NewEmployeeGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("Add", mock.Anything).Return(tt.ucResult, tt.ucErr)
			}

			jsonBytes, _ := json.Marshal(tt.args.req)
			reqBody := bytes.NewBuffer(jsonBytes)

			req, errReq := http.NewRequest(http.MethodPost, "/employees", reqBody)
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

func TestEmployeeGinHandler_Detail(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		ucResult   *model.Employee
		ucErr      error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{id: "1"},
			ucResult: &model.Employee{
				ID:      1,
				Name:    "test",
				Gender:  "laki-laki",
				Grade:   1,
				Married: false,
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
			mockUC := mocks.NewEmployeeUsecaseI(t)

			h := handler.NewEmployeeGinHandler(mockUC)

			if tt.ucResult != nil || tt.ucErr != nil {
				mockUC.On("Detail", mock.Anything).Return(tt.ucResult, tt.ucErr)
			}

			req, errReq := http.NewRequest(http.MethodGet, "/employees/"+tt.args.id, nil)
			assert.NoError(t, errReq)

			// res := httptest.NewRecorder()
			// c, _ := gin.CreateTestContext(res)
			// c.Request = req

			// h.Detail(c)

			res := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/employees/:id", h.Detail)
			router.ServeHTTP(res, req)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
		})
	}
}
