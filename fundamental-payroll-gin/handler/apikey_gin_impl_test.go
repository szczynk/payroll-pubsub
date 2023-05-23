package handler_test

import (
	"bytes"
	"encoding/json"
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeyGinHandler_Generate(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name:       "success",
			args:       args{name: "bagus"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid field",
			args:       args{name: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong url",
			args:       args{name: "bagus"},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := handler.NewAPIKeyGinHandler("http://localhost:8081")
			if strings.Contains(tt.name, "wrong url") {
				h = handler.NewAPIKeyGinHandler("http://localhost:8082")
			}

			apikey := model.APIKeyReq{
				Name: tt.args.name,
			}
			jsonBytes, _ := json.Marshal(apikey)
			reqBody := bytes.NewBuffer(jsonBytes)

			req, errReq := http.NewRequest(http.MethodPost, "/generate", reqBody)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			c.Request = req

			h.Generate(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code should be OK")
		})
	}
}
