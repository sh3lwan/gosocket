package models

type Response struct {
	StatusCode int16 `json:"status_code"`
	Data       map[string]any
}

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
