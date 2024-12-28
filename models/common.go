package models

// The generic successful operation response.
//
// @description	Operation Succeeded.
type SuccessResp struct {
	Success bool `json:"success"`
}

// This is the error struct returned from all requests and
// others are just used for documentation purposes.
//
// @description	Operation Failed.
// @description	Maybe invalid payload or internal server error.
type ErrResp struct {
	Message any `json:"message" example:"Malformed payload"` // Error cause.
	Code    int `json:"code" example:"400"`                  // HTTP code.
}

// Used, for e.g, for database related error.
//
// @description	Something wrong happened on the server.
type InternalServerError struct {
	Message any `json:"message" example:"Something went wrong"` // Error cause.
	Code    int `json:"code" example:"500"`                     // HTTP code.
}

// BadRequestError used for errors of validation.
//
// @description	Request payload failed validation.
type BadRequestError struct {
	Message any `json:"message" example:"Payload Failed Validation"` // Error cause.
	Code    int `json:"code" example:"422"`                          // HTTP code.
}

// UnprocessableEntityError used for errors of parsing request body.
//
// @description	Failed to parse payload.
type UnprocessableEntityError struct {
	Message any `json:"message" example:"Payload Parsing Failed"` // Error cause.
	Code    int `json:"code" example:"422"`                       // HTTP code.
}
