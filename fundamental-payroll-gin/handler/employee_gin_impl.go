package handler

import (
	"database/sql"
	"errors"
	"fundamental-payroll-gin/helper/response"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeGinHandler struct {
	employeeUC usecase.EmployeeUsecaseI
}

func NewEmployeeGinHandler(employeeUC usecase.EmployeeUsecaseI) EmployeeGinHandlerI {
	h := new(EmployeeGinHandler)
	h.employeeUC = employeeUC
	return h
}

func (h *EmployeeGinHandler) List(c *gin.Context) {
	employees, ucErr := h.employeeUC.List()
	if ucErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, employees)
}

func (h *EmployeeGinHandler) Add(c *gin.Context) {
	employeeReq := new(model.EmployeeRequest)
	if bindErr := c.ShouldBindJSON(&employeeReq); bindErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, bindErr.Error())
		return
	}

	if employeeReq.Gender != "laki-laki" && employeeReq.Gender != "perempuan" {
		response.NewJSONResErr(c, http.StatusBadRequest, `invalid gender, try "laki-laki" or "perempuan"`)
		return
	}

	if employeeReq.Grade <= 0 {
		response.NewJSONResErr(c, http.StatusBadRequest, "invalid grade")
		return
	}

	employee, ucErr := h.employeeUC.Add(employeeReq)
	if ucErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, employee)
}

func (h *EmployeeGinHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, err.Error())
		return
	}

	employee, ucErr := h.employeeUC.Detail(id)
	if ucErr != nil {
		if errors.Is(ucErr, sql.ErrNoRows) {
			response.NewJSONResErr(c, http.StatusNotFound, "employee with that id not found")
			return
		}
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, employee)
}
