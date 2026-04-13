package utils

import "github.com/gin-gonic/gin"

// RespondError writes a standard error envelope.
func RespondError(c *gin.Context, status int, mensagem string) {
	c.JSON(status, gin.H{"mensagem": mensagem, "data": nil})
}

// RespondOK writes a standard success envelope.
func RespondOK(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"mensagem": "ok", "data": data})
}
