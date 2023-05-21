package model

type APIKey struct {
	Key string `json:"api_key"`
}

type APIKeyReq struct {
	Name string `form:"name" json:"name" binding:"required"`
}
