package response

import "github.com/labstack/echo/v5"

type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func Send(c *echo.Context, data interface{}, message string, status int) error {
	return c.JSON(status, APIResponse{
		Data:    data,
		Message: message,
		Status:  status,
	})
}
