package models

// @description Operation Succeeded.
type SuccessResp struct {
	Success bool `json:"success"`
}

// @description Operation Failed.
// @description Maybe invalid payload or internal server error.
type ErrResp struct {
	Message any `json:"message" example:"Malformed payload"` // Error cause.
	Code    int `json:"code" example:"400"`                  // HTTP code.
}

// @description Something wrong happened on the server.
type InternalServerError struct {
	Message any `json:"message" example:"Something went wrong"` // Error cause.
	Code    int `json:"code" example:"500"`                     // HTTP code.
}
