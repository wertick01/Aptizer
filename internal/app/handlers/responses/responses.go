package responses

import (
	"github.com/gin-gonic/gin"
)

// responses.WrapGinErrorWithStatus - Forms a response in case of an error.
func WrapGinErrorWithStatus(c *gin.Context, err, clienterr error, httpStatus int) {
	c.AbortWithStatusJSON(httpStatus, gin.H{
		"result":  "error",
		"code":    httpStatus,
		"message": clienterr.Error(),
		"detail":  err.Error(),
	})
}
