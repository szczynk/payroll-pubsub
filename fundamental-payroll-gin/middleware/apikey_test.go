package middleware_test

import (
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/helper/logger"
	"fundamental-payroll-gin/middleware"
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPIKey(t *testing.T) {
	type args struct {
		id     string
		apiKey string
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
				id:     "1",
				apiKey: "uL7buMVaRdhd/FGX5pFS0wB4ojyHTNRD2XAKx3pZOmBCK63URNHLZak=",
			},
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
			name: "missing api key",
			args: args{
				id:     "1",
				apiKey: "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid api key",
			args: args{
				id:     "1",
				apiKey: "shit",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong url",
			args: args{
				id: "1",
			},
			wantStatus: http.StatusUnauthorized,
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
			req.Header.Add("X-API-Key", tt.args.apiKey)

			res := httptest.NewRecorder()
			router := gin.Default()

			logger := logger.New(true)
			router.Use(middleware.Logger(logger))
			router.Use(middleware.APIKey("http://localhost:8081"))
			if strings.Contains(tt.name, "wrong url") {
				router.Use(middleware.APIKey("http://localhost:8082"))
			}
			router.GET("/employees/:id", h.Detail)
			router.ServeHTTP(res, req)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
		})
	}
}
