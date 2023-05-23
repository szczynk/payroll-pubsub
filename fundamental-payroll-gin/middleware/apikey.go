package middleware

import (
	"encoding/json"
	"fundamental-payroll-gin/helper/response"
	"fundamental-payroll-gin/helper/timeout"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKey(apiVerificationURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		ctx, cancel := timeout.NewCtxTimeout()
		defer cancel()

		req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, apiVerificationURL+"/verify", nil)
		if errReq != nil {
			c.Abort()
			response.NewJSONResErr(c, http.StatusInternalServerError, errReq.Error())
			return
		}
		req.Header.Add("X-API-Key", apiKey)

		res, errRes := http.DefaultClient.Do(req)
		if errRes != nil {
			c.Abort()
			response.NewJSONResErr(c, http.StatusInternalServerError, errRes.Error())
			return
		}
		defer res.Body.Close()

		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			c.Abort()
			response.NewJSONResErr(c, http.StatusInternalServerError, readErr.Error())
			return
		}

		var jsonRes response.JSONRes
		if jsonErr := json.Unmarshal(body, &jsonRes); jsonErr != nil {
			c.Abort()
			response.NewJSONResErr(c, http.StatusInternalServerError, jsonErr.Error())
			return
		}

		statusCode := res.StatusCode
		if statusCode >= http.StatusBadRequest && statusCode <= http.StatusInternalServerError {
			c.Abort()
			response.NewJSONResErr(c, jsonRes.Status, jsonRes.Error)
			return
		}

		c.Next()
	}
}
