package models

type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
