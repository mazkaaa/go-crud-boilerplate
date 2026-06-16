package response

import "github.com/labstack/echo/v5"

type APIResponse struct {
	Data    interface{}     `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Message string          `json:"message"`
	Status  int             `json:"status"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func Send(c *echo.Context, data interface{}, message string, status int) error {
	return c.JSON(status, APIResponse{
		Data:    data,
		Message: message,
		Status:  status,
	})
}

func SendPaginated(c *echo.Context, data interface{}, meta PaginationMeta, message string, status int) error {
	return c.JSON(status, APIResponse{
		Data:    data,
		Meta:    &meta,
		Message: message,
		Status:  status,
	})
}
