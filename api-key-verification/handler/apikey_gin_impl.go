package handler

import (
	"api-key-verification/helper/crypt"
	"api-key-verification/helper/response"
	"api-key-verification/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIKeyGinHandler struct {
	passphrase string
}

func NewAPIKeyGinHandler(passphrase string) APIKeyGinHandlerI {
	h := new(APIKeyGinHandler)
	h.passphrase = passphrase
	return h
}

func (h *APIKeyGinHandler) Generate(c *gin.Context) {
	apiKeyReq := new(model.APIKeyReq)
	if bindErr := c.ShouldBindJSON(apiKeyReq); bindErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, bindErr.Error())
		return
	}

	apiKey, encryptErr := crypt.Encrypt(apiKeyReq.Name+"|phincon", h.passphrase)
	if encryptErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, encryptErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, model.APIKey{Key: apiKey})
}
