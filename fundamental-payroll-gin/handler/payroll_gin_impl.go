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

type PayrollGinHandler struct {
	payrollUC usecase.PayrollUsecaseI
}

func NewPayrollGinHandler(payrollUC usecase.PayrollUsecaseI) PayrollGinHandlerI {
	h := new(PayrollGinHandler)
	h.payrollUC = payrollUC
	return h
}

func (h *PayrollGinHandler) List(c *gin.Context) {
	payrolls, ucErr := h.payrollUC.List()
	if ucErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, payrolls)
}

func (h *PayrollGinHandler) Add(c *gin.Context) {
	payrollReq := new(model.PayrollRequest)
	if bindErr := c.ShouldBindJSON(payrollReq); bindErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, bindErr.Error())
		return
	}

	if payrollReq.EmployeeID <= 0 {
		response.NewJSONResErr(c, http.StatusBadRequest, "invalid employee's id")
		return
	}

	if payrollReq.TotalHariMasuk < 0 {
		response.NewJSONResErr(c, http.StatusBadRequest, "invalid total hari masuk")
		return
	}
	if payrollReq.TotalHariTidakMasuk < 0 {
		response.NewJSONResErr(c, http.StatusBadRequest, "invalid total hari tidak masuk")
		return
	}

	payroll, ucErr := h.payrollUC.Add(payrollReq)
	if ucErr != nil {
		if errors.Is(ucErr, sql.ErrNoRows) {
			response.NewJSONResErr(c, http.StatusBadRequest, "invalid employee's id")
			return
		}
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, payroll)
}

func (h *PayrollGinHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, err.Error())
		return
	}

	payroll, ucErr := h.payrollUC.Detail(id)
	if ucErr != nil {
		if errors.Is(ucErr, sql.ErrNoRows) {
			response.NewJSONResErr(c, http.StatusNotFound, "payroll with that id not found")
			return
		}
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, payroll)
}
