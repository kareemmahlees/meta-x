package models

type SuccessResp struct {
	Success bool `json:"success"`
}

type ErrResp struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}
