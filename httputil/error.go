package httputil

import "github.com/gin-gonic/gin"

// NewError .
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError .
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
