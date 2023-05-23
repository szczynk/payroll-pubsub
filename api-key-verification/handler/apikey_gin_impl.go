package handler

import (
	"api-key-verification/helper/crypt"
	"api-key-verification/helper/response"
	"api-key-verification/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIKeyGinHandler struct {
	crypt crypt.Crypt
}

func NewAPIKeyGinHandler(crypt crypt.Crypt) APIKeyGinHandlerI {
	h := new(APIKeyGinHandler)
	h.crypt = crypt
	return h
}

func (h *APIKeyGinHandler) Generate(c *gin.Context) {
	apiKeyReq := new(model.APIKeyReq)
	if bindErr := c.ShouldBindJSON(&apiKeyReq); bindErr != nil {
		response.NewJSONResErr(c, http.StatusBadRequest, bindErr.Error())
		return
	}

	apiKey, encryptErr := h.crypt.Encrypt(apiKeyReq.Name + "|phincon")
	if encryptErr != nil {
		response.NewJSONResErr(c, http.StatusInternalServerError, encryptErr.Error())
		return
	}

	response.NewJSONRes(c, http.StatusOK, model.APIKey{Key: apiKey})
}

func (h *APIKeyGinHandler) Verify(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		response.NewJSONResErr(c, http.StatusUnauthorized, "API key is missing")
		return
	}

	// Verify the API key here using AES encryption
	// passphrase := "your-secret-passphrase" // This should be a secret stored securely
	decryptedStr, errDecrypt := h.crypt.Decrypt(apiKey)
	if errDecrypt != nil && !strings.Contains(decryptedStr, "phincon") { // if error and didn't contain "phincon"
		response.NewJSONResErr(c, http.StatusUnauthorized, "Invalid API key")
		return
	}

	response.NewJSONRes(c, http.StatusOK, nil)
}
