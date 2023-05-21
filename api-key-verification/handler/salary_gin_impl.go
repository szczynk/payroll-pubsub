package handler

import (
	"api-key-verification/helper/response"
	"api-key-verification/helper/timeout"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalaryGinHandler struct {
	apiURL string
}

func NewSalaryGinHandler(apiURL string) SalaryGinHandlerI {
	h := new(SalaryGinHandler)
	h.apiURL = apiURL
	return h
}

func (h *SalaryGinHandler) List(c *gin.Context) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, h.apiURL+"/salaries", nil)
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
	if jsonErr := json.Unmarshal(body, &jsonRes); jsonErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, jsonErr.Error())
		return
	}

	c.JSON(http.StatusOK, jsonRes)
}
