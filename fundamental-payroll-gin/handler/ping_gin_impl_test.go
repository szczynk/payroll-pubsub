package handler_test

import (
	"fundamental-payroll-gin/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingGinHandler_Ping(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		// TODO: Add test cases.
		{
			name:       "success",
			wantStatus: http.StatusOK,
			wantBody:   "pong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := handler.NewPingGinHandler()

			req, errReq := http.NewRequest(http.MethodGet, "/ping", nil)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			c.Request = req

			h.Ping(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			// Use assert package to simplify assertions
			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
			assert.Contains(t, res.Body.String(), tt.wantBody, "response body should be contains 'pong'")
		})
	}
}
