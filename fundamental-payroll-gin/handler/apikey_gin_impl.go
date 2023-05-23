package handler

import (
	"bytes"
	"encoding/json"
	"fundamental-payroll-gin/helper/response"
	"fundamental-payroll-gin/helper/timeout"
	"fundamental-payroll-gin/model"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIKeyGinHandler struct {
	apiKeyVerificationURL string
}

func NewAPIKeyGinHandler(apiKeyVerificationURL string) APIKeyGinHandlerI {
	h := new(APIKeyGinHandler)
	h.apiKeyVerificationURL = apiKeyVerificationURL
	return h
}

func (h *APIKeyGinHandler) Generate(c *gin.Context) {
	apiKeyReq := new(model.APIKeyReq)
	if bindErr := c.ShouldBindJSON(&apiKeyReq); bindErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, bindErr.Error())
		return
	}

	apiKeyReqBytes, jsonErr := json.Marshal(apiKeyReq)
	if jsonErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, jsonErr.Error())
		return
	}

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, h.apiKeyVerificationURL+"/generate",
		bytes.NewBuffer(apiKeyReqBytes))
	if reqErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, reqErr.Error())
		return
	}

	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, resErr.Error())
		return
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, readErr.Error())
		return
	}

	var jsonRes response.JSONRes
	if jsonUErr := json.Unmarshal(body, &jsonRes); jsonUErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, jsonUErr.Error())
		return
	}

	c.JSON(http.StatusOK, jsonRes)
}
