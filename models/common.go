package models

type SuccessResp struct {
	Success bool `json:"success"`
}

type ErrResp struct {
	Message any `json:"message"`
	Code    int `json:"code"`
}
