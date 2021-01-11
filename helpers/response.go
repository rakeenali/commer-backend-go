package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	statusOk      = "OK"
	statusErr     = "ERROR"
	statusBodyErr = "BODY_ERROR"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
}

// InvalidBodyResponse will return an invalid if request body is not valid
func InvalidBodyErrorResponse(c *gin.Context, data interface{}) {
	var res Response

	res.Message = "Invalid body"
	res.Data = data
	res.Status = statusBodyErr
	c.JSON(http.StatusBadRequest, res)
}

// InternalServerErrorResponse will return public message
func InternalServerErrorResponse(c *gin.Context, err error) {
	var res Response

	fmt.Println("InternalServerErrorResponse", err)

	res.Message = "Something went wrong"
	res.Data = nil
	res.Status = statusErr
	c.JSON(http.StatusInternalServerError, res)
}

// OKResponse will return a successfull response
func OKResponse(c *gin.Context, message string, code int, data interface{}) {
	var res Response

	httpCode := code
	res.Message = message

	if message == "" {
		res.Message = "Success"
	}
	if code == 0 {
		httpCode = http.StatusOK
	}

	res.Data = data
	res.Status = statusOk
	c.JSON(httpCode, res)
}

// ErrResponse will return a response with error
func ErrResponse(c *gin.Context, message string, code int, data interface{}) {
	var res Response

	httpCode := code
	res.Message = message

	if message == "" {
		res.Message = "Error"
	}
	if code == 0 {
		httpCode = http.StatusNotFound
	}

	res.Data = data
	res.Status = statusErr
	c.JSON(httpCode, res)
}
