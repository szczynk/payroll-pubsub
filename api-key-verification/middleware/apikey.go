package middleware

import (
	"api-key-verification/helper/crypt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKey(passphrase string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is missing"})
			return
		}

		// Verify the API key here using AES encryption
		// passphrase := "your-secret-passphrase" // This should be a secret stored securely
		decryptedStr, err := crypt.Decrypt(apiKey, passphrase)
		if err != nil && !strings.Contains(decryptedStr, "phincon") { // if error and didn't contain "phincon"
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Next()
	}
}
