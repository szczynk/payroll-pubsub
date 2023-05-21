//nolint:dupl // this is why
package handler

import (
	"api-key-verification/helper/response"
	"api-key-verification/helper/timeout"
	"api-key-verification/model"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeeGinHandler struct {
	apiURL string
}

func NewEmployeeGinHandler(apiURL string) EmployeeGinHandlerI {
	h := new(EmployeeGinHandler)
	h.apiURL = apiURL
	return h
}

func (h *EmployeeGinHandler) List(c *gin.Context) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, h.apiURL+"/employees", nil)
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

func (h *EmployeeGinHandler) Add(c *gin.Context) {
	employeeReq := new(model.EmployeeRequest)
	if bindErr := c.ShouldBindJSON(employeeReq); bindErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, bindErr.Error())
		return
	}

	employeeReqBytes, jsonErr := json.Marshal(employeeReq)
	if jsonErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, jsonErr.Error())
		return
	}

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, h.apiURL+"/employees",
		bytes.NewBuffer(employeeReqBytes))
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

func (h *EmployeeGinHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, h.apiURL+"/employees/"+id, nil)
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
