package global

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/shared/utils"
)

type errorMapping struct {
	target error
	status int
}

func respondMappedError(c *gin.Context, err error, mappings ...errorMapping) {
	for _, mapping := range mappings {
		if errors.Is(err, mapping.target) {
			utils.RespondError(c, mapping.status, err.Error())
			return
		}
	}

	utils.RespondError(c, http.StatusInternalServerError, err.Error())
}
