package handler

import (
	"fundamental-payroll-gin/helper/response"
	"fundamental-payroll-gin/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalaryGinHandler struct {
	salaryUC usecase.SalaryUsecaseI
}

func NewSalaryGinHandler(salaryUC usecase.SalaryUsecaseI) SalaryGinHandlerI {
	h := new(SalaryGinHandler)
	h.salaryUC = salaryUC
	return h
}

func (h *SalaryGinHandler) List(c *gin.Context) {
	salaries, ucErr := h.salaryUC.List()
	if ucErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, ucErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, salaries)
}
